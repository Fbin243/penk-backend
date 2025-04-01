import { readFileSync } from "fs";
import OpenAI from "openai";
import { Uploadable } from "openai/core.mjs";
import {
  ChatCompletionAssistantMessageParam,
  ChatCompletionMessageParam,
  ChatCompletionUserMessageParam,
  CompletionUsage,
} from "openai/resources/index.mjs";
import { encoding_for_model } from "tiktoken";

import { Message, MessageType } from "../../utils/types";

const client = new OpenAI({
  apiKey: process.env.OPEN_AI_API_KEY,
});

const instructions = readFileSync("instructions.md", "utf8");

const analyzeUsage = (usage: CompletionUsage) => {
  /*
    https://platform.openai.com/docs/pricing
    Pricing of gpt-4o-mini-audio-preview:
    - Text Tokens Input: $0.150 / 1M tokens
    - Cached Text Tokens Input: $0.075 / 1M tokens
    - Text Tokens Output: $0.600 / 1M tokens

    - Audio Tokens Input: $10.000 / 1M tokens
    - Audio Tokens Output: $20.000 / 1M tokens
  */

  // @ts-expect-error - OpenAI types are not updated
  const textInputTokens = usage.prompt_tokens_details?.text_tokens || 0;
  const cachedTextTokens = usage.prompt_tokens_details?.cached_tokens || 0;
  const audioInputTokens = usage.prompt_tokens_details?.audio_tokens || 0;

  // @ts-expect-error - OpenAI types are not updated
  const textOutputTokens = usage.completion_tokens_details?.text_tokens || 0;
  const audioOutputTokens = usage.completion_tokens_details?.audio_tokens || 0;

  const textInputCost = (textInputTokens * 0.15 + cachedTextTokens * 0.075) / 1000000;
  const audioInputCost = audioInputTokens / 1000000;
  const textOutputCost = (textOutputTokens * 0.6) / 1000000;
  const audioOutputCost = (audioOutputTokens * 20) / 1000000;

  const totalCost = textInputCost + audioInputCost + textOutputCost + audioOutputCost;

  const usdToVnd = 25575;

  return {
    textInputTokens,
    cachedTextTokens,
    audioInputTokens,
    textOutputTokens,
    audioOutputTokens,
    textInputCost,
    audioInputCost,
    textOutputCost,
    audioOutputCost,
    totalCost,
    totalCostVnd: `${Math.round(
      totalCost * usdToVnd,
    ).toLocaleString()} VND, with 1 USD = ${usdToVnd} VND`,
  };
};

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

  const encoder = encoding_for_model("gpt-4o");
  const tokens = encoder.encode(transcription.text).length;
  // Pricing of gpt-4o-transcribe: $6.00 per 1M tokens
  const transcriptionCost = (tokens * 6) / 1000000;

  return {
    text: transcription.text,
    cost: transcriptionCost,
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
  onComplete?: (fullMessage: Message, usage?: ReturnType<typeof analyzeUsage>) => void,
) => {
  const openAiMessages: ChatCompletionMessageParam[] = [
    {
      role: "system",
      content: instructions,
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
      max_tokens: 2048,
    });

    const aiTimestamp = new Date().toISOString();

    for await (const chunk of stream) {
      const content = chunk.choices[0]?.delta?.content || "";
      if (content) {
        fullContent += content;
        onChunk(content, aiTimestamp);
      }
    }

    // Create the complete message object after streaming is done
    const aiMessage: Message = {
      type: MessageType.AiMessage,
      content: fullContent,
      timestamp: aiTimestamp,
    };

    // Call onComplete callback if provided
    if (onComplete) {
      onComplete(aiMessage, undefined); // We don't have usage info in streaming
    }

    return aiMessage;
  } catch (error) {
    console.error("Error in textChatStream:", error);
    throw error;
  }
};

export const chat = async (props: {
  userData: string;
  history: Message[];
  newMessage: string;
  voiceMode?: boolean;
}) => {
  const openAiMessages: ChatCompletionMessageParam[] = [
    {
      role: "system",
      content: instructions,
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

  const completion = await client.chat.completions.create({
    model: props.voiceMode ? "gpt-4o-mini-audio-preview" : "gpt-4o-mini",
    messages: openAiMessages,
    modalities: props.voiceMode ? ["text", "audio"] : ["text"],
    audio: props.voiceMode
      ? {
          format: "wav",
          voice: "alloy",
        }
      : undefined,
    max_completion_tokens: 2048,
  });

  const aiMessage: Message = {
    type: MessageType.AiMessage,
    content:
      completion.choices[0].message.content ||
      completion.choices[0].message.audio?.transcript ||
      "",
    timestamp: new Date().toISOString(),
  };

  return {
    aiMessage,
    audio: completion.choices[0].message.audio,
    usage: completion.usage ? analyzeUsage(completion.usage) : undefined,
  };
};
