import { readFileSync } from "fs";
import OpenAI from "openai";
import { Uploadable } from "openai/core.mjs";
import {
  ChatCompletionAssistantMessageParam,
  ChatCompletionMessageParam,
  ChatCompletionToolMessageParam,
  ChatCompletionUserMessageParam,
} from "openai/resources/index.mjs";

import { getAudioDuration } from "../audio";
import { Message, MessageType } from "../types";
import { callFunction, FunctionName, tools } from "./functions";
import {
  calculateCompletionUsage,
  calculateTranscriptionCost,
  calculateTtsCost,
  gpt4oMiniPricingModel,
} from "./pricing";

const client = new OpenAI({
  apiKey: process.env.OPEN_AI_API_KEY,
});

const baseInstruction = readFileSync("resources/instructions/base.md", "utf8");
const voiceModeInstruction = readFileSync("resources/instructions/voice-mode.md", "utf8");

export const base64ToUploadable = (
  base64: string,
  filename: string,
  mimeType: string = "audio/mpeg",
): File => {
  // Remove the Base64 prefix if it exists
  const base64Data = base64.split(",")[1] || base64;

  // Convert Base64 to binary
  const byteCharacters = atob(base64Data);
  const byteNumbers = new Array(byteCharacters.length);
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }
  const byteArray = new Uint8Array(byteNumbers);

  // Create a Blob, then a File
  return new File([byteArray], filename, { type: mimeType });
};

export const transcribeAudio = async (audio: Uploadable) => {
  const transcription = await client.audio.transcriptions.create({
    file: audio,
    model: "gpt-4o-transcribe",
  });

  return {
    text: transcription.text,
    cost: calculateTranscriptionCost(transcription.text),
  };
};

/**
 * Streams chat completions from OpenAI with text-only responses
 * @param props Chat properties including user data, history, and the new message
 * @param onChunk Callback function that receives each chunk of the streamed response
 * @param onComplete Optional callback function called when streaming is complete with the full message
 */
export const textChatStream = async (
  props: {
    userData: string;
    history: Message[];
    newMessage: string;
  },
  onChunk: (chunk: string, timestamp: string) => void,
  onComplete?: (fullMessage: Message) => void,
) => {
  const openAiMessages: ChatCompletionMessageParam[] = [
    {
      role: "system",
      content: baseInstruction,
    },
    {
      role: "user",
      content: `My data: ${props.userData}`,
    },
    ...props.history.map((message) => {
      if (message.type === MessageType.UserMessage) {
        return {
          role: "user",
          content: message.content,
        } satisfies ChatCompletionUserMessageParam;
      }
      return {
        role: "assistant",
        content: message.content,
      } satisfies ChatCompletionAssistantMessageParam;
    }),
    {
      role: "user",
      content: props.newMessage,
    },
  ];

  let fullContent = "";

  try {
    const stream = await client.chat.completions.create({
      model: "gpt-4o-mini",
      messages: openAiMessages,
      stream: true,
      stream_options: {
        include_usage: true,
      },
      max_tokens: 2048,
      tools,
      tool_choice: "auto",
    });

    const aiTimestamp = new Date().toISOString();

    let cost = 0;

    const finalToolCalls: Record<
      number,
      OpenAI.Chat.Completions.ChatCompletionChunk.Choice.Delta.ToolCall
    > = {};

    for await (const chunk of stream) {
      const toolCalls = chunk.choices[0]?.delta?.tool_calls || [];
      if (toolCalls.length > 0) {
        for (const toolCall of toolCalls) {
          const { index } = toolCall;

          if (!finalToolCalls[index]) {
            finalToolCalls[index] = toolCall;
          }

          if (finalToolCalls[index].function && toolCall.function?.arguments) {
            finalToolCalls[index].function.arguments += toolCall.function.arguments;
          }
        }
      }

      const content = chunk.choices[0]?.delta?.content || "";
      if (content) {
        fullContent += content;
        onChunk(content, aiTimestamp);
      }

      if (chunk.usage) {
        cost += calculateCompletionUsage(chunk.usage, gpt4oMiniPricingModel);
      }
    }

    if (Object.keys(finalToolCalls).length > 0) {
      for (const toolCall of Object.values(finalToolCalls)) {
        const result = await callFunction(
          toolCall.function!.name as FunctionName,
          JSON.parse(toolCall.function!.arguments!),
        );
        openAiMessages.push({
          role: "assistant",
          tool_calls: [
            {
              id: toolCall.id!,
              function: {
                name: toolCall.function!.name!,
                arguments: toolCall.function!.arguments!,
              },
              type: "function",
            },
          ],
        });
        openAiMessages.push({
          role: "tool",
          content: JSON.stringify(result),
          tool_call_id: toolCall.id!,
        } satisfies ChatCompletionToolMessageParam);
      }

      const completion = await client.chat.completions.create({
        model: "gpt-4o-mini",
        messages: openAiMessages,
        max_tokens: 2048,
        tools,
        tool_choice: "auto",
        stream: true,
      });

      for await (const chunk of completion) {
        const content = chunk.choices[0]?.delta?.content || "";
        if (content) {
          fullContent += content;
          onChunk(content, aiTimestamp);
        }

        if (chunk.usage) {
          cost += calculateCompletionUsage(chunk.usage, gpt4oMiniPricingModel);
        }
      }
    }

    const aiMessage: Message = {
      type: MessageType.AiMessage,
      content: fullContent || "Something went wrong",
      timestamp: aiTimestamp,
    };

    if (onComplete) {
      onComplete(aiMessage);
    }

    return {
      aiMessage,
      cost,
    };
  } catch (error) {
    console.error("Error in textChatStream:", error);
    throw error;
  }
};

/**
 * Process audio chat with OpenAI, including audio response generation
 * @param props Chat properties including user data, history, and the transcribed message
 * @param onTextReady Callback function for when the text response is ready
 * @param onAudioReady Callback function for when the audio response is ready
 */
export const audioChat = async (
  props: {
    userData: string;
    history: Message[];
    transcribedMessage: string;
  },
  onTextReady: (textResponse: Message) => void,
  onAudioReady: (base64Audio: string) => void,
) => {
  const openAiMessages: ChatCompletionMessageParam[] = [
    {
      role: "system",
      content: baseInstruction,
    },
    {
      role: "user",
      content: `My data: ${props.userData}`,
    },
    ...props.history.map((message) => {
      if (message.type === MessageType.UserMessage) {
        return {
          role: "user",
          content: message.content,
        } satisfies ChatCompletionUserMessageParam;
      }
      return {
        role: "assistant",
        content: message.content,
      } satisfies ChatCompletionAssistantMessageParam;
    }),
    {
      role: "user",
      content: props.transcribedMessage,
    },
  ];

  try {
    const completion = await client.chat.completions.create({
      model: "gpt-4o-mini",
      messages: openAiMessages,
      modalities: ["text"],
      max_completion_tokens: 2048,
    });

    const aiMessage: Message = {
      type: MessageType.AiMessage,
      content: completion.choices[0].message.content || "",
      timestamp: new Date().toISOString(),
    };

    onTextReady(aiMessage);
    console.log("Text response sent:", aiMessage.content);

    let cost = completion.usage
      ? calculateCompletionUsage(completion.usage, gpt4oMiniPricingModel)
      : 0;

    const mp3 = await client.audio.speech.create({
      model: "gpt-4o-mini-tts",
      voice: "ash",
      input: aiMessage.content,
      instructions: voiceModeInstruction,
    });
    const mp3ArrayBuffer = await mp3.arrayBuffer();
    const base64Audio = Buffer.from(mp3ArrayBuffer).toString("base64");
    onAudioReady(base64Audio);

    const audioDuration = await getAudioDuration(Buffer.from(mp3ArrayBuffer));
    const ttsCost = calculateTtsCost(audioDuration);
    console.log(`TTS cost for ${audioDuration} seconds: ${ttsCost}`);
    cost += ttsCost;

    return {
      aiMessage,
      audio: base64Audio,
      cost,
    };
  } catch (error) {
    console.error("Error in audioChatStream:", error);
    throw error;
  }
};
