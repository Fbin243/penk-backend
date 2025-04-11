import { readFileSync } from "fs";
import OpenAI from "openai";
import { Uploadable } from "openai/core.mjs";
import {
  ChatCompletionMessageParam,
  ChatCompletionToolMessageParam,
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
import { setupInitialMessages, streamAssistantResponse } from "./utils";

const client = new OpenAI({
  apiKey: process.env.OPEN_AI_API_KEY,
});

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

export const textChatStream = async (
  props: {
    userData: string;
    history: Message[];
    newMessage: string;
  },
  onChunk: (chunk: string, timestamp: string) => void,
  onComplete?: (fullMessage: Message) => void,
) => {
  const aiTimestamp = new Date().toISOString();
  let fullContent = "";
  let cost = 0;

  const openAiMessages = setupInitialMessages(props.userData, props.history);
  openAiMessages.push({ role: "user", content: props.newMessage });

  try {
    const { content, toolCalls, usage } = await streamAssistantResponse({
      client,
      messages: openAiMessages,
      aiTimestamp,
      onChunk,
    });
    fullContent += content;
    if (usage) cost += calculateCompletionUsage(usage, gpt4oMiniPricingModel);

    let currentToolCalls = toolCalls;

    while (Object.keys(currentToolCalls).length > 0) {
      const newMessages: ChatCompletionMessageParam[] = [];

      for (const toolCall of Object.values(currentToolCalls)) {
        const args = JSON.parse(toolCall.function!.arguments!);

        let result: unknown;
        try {
          result = await callFunction(toolCall.function!.name as FunctionName, args);
        } catch (error) {
          result = error.message || error;
          console.error("Error calling function:", result);
        }

        newMessages.push({
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

        newMessages.push({
          role: "tool",
          content: JSON.stringify(result),
          tool_call_id: toolCall.id!,
        } satisfies ChatCompletionToolMessageParam);
      }

      openAiMessages.push(...newMessages);

      const followup = await streamAssistantResponse({
        client,
        messages: openAiMessages,
        aiTimestamp,
        onChunk,
      });
      fullContent += followup.content;
      if (followup.usage) cost += calculateCompletionUsage(followup.usage, gpt4oMiniPricingModel);
      currentToolCalls = followup.toolCalls;
    }

    const aiMessage: Message = {
      type: MessageType.AiMessage,
      content: fullContent || "Something went wrong",
      timestamp: aiTimestamp,
    };

    if (onComplete) onComplete(aiMessage);
    return { aiMessage, cost };
  } catch (error) {
    console.error("Error in textChatStream:", error);
    throw error;
  }
};

export const audioChat = async (
  props: {
    userData: string;
    history: Message[];
    transcribedMessage: string;
  },
  onTextReady: (textResponse: Message) => void,
  onAudioReady: (base64Audio: string) => void,
) => {
  const openAiMessages = setupInitialMessages(props.userData, props.history);
  openAiMessages.push({
    role: "user",
    content: props.transcribedMessage,
  });

  let finalContent = "";
  const aiTimestamp = new Date().toISOString();
  let totalCost = 0;

  while (true) {
    const completion = await client.chat.completions.create({
      model: "gpt-4o-mini",
      messages: openAiMessages,
      modalities: ["text"],
      tools,
      tool_choice: "auto",
      max_completion_tokens: 2048,
    });

    const choice = completion.choices[0].message;
    totalCost += completion.usage
      ? calculateCompletionUsage(completion.usage, gpt4oMiniPricingModel)
      : 0;

    // If tool calls exist, call them and loop again
    if (choice.tool_calls?.length) {
      for (const toolCall of choice.tool_calls) {
        const fnName = toolCall.function.name as FunctionName;
        const args = JSON.parse(toolCall.function.arguments || "{}");

        let result: unknown;
        try {
          result = await callFunction(fnName, args);
        } catch (error) {
          result = error.message || error;
          console.error("Error calling function:", result);
        }

        openAiMessages.push({
          role: "assistant",
          tool_calls: [toolCall],
        });

        openAiMessages.push({
          role: "tool",
          tool_call_id: toolCall.id,
          content: JSON.stringify(result),
        });
      }
      continue; // back to top with new messages
    }

    // Final message
    finalContent = choice.content || "";
    break;
  }

  const aiMessage: Message = {
    type: MessageType.AiMessage,
    content: finalContent,
    timestamp: aiTimestamp,
  };

  onTextReady(aiMessage);
  console.log("Text response sent:", finalContent);

  // Generate audio
  const mp3 = await client.audio.speech.create({
    model: "gpt-4o-mini-tts",
    voice: "ash",
    input: aiMessage.content,
    instructions: voiceModeInstruction,
  });

  const mp3ArrayBuffer = await mp3.arrayBuffer();
  const base64Audio = Buffer.from(mp3ArrayBuffer).toString("base64");
  onAudioReady(base64Audio);

  const duration = await getAudioDuration(Buffer.from(mp3ArrayBuffer));
  totalCost += calculateTtsCost(duration);
  console.log(`TTS cost for ${duration} seconds: ${calculateTtsCost(duration)}`);

  return {
    aiMessage,
    audio: base64Audio,
    cost: totalCost,
  };
};
