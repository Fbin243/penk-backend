import ffmpeg from "fluent-ffmpeg";
import { promises as fs, readFileSync } from "fs";
import OpenAI from "openai";
import { Uploadable } from "openai/core.mjs";
import {
  ChatCompletionAssistantMessageParam,
  ChatCompletionMessageParam,
  ChatCompletionUserMessageParam,
  CompletionUsage,
} from "openai/resources/index.mjs";
import os from "os";
import path from "path";
import { encoding_for_model } from "tiktoken";

import { Message, MessageType } from "../types";

const client = new OpenAI({
  apiKey: process.env.OPEN_AI_API_KEY,
});

const baseInstruction = readFileSync("resources/instructions/base.md", "utf8");

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
  };
};

/**
 * Converts a File object from one audio format to another
 * @param audioFile The input audio file
 * @returns A Promise that resolves to a new File in the target format
 */
export const convertAudioFormatToMp3 = async (audioFile: File): Promise<File> => {
  // Create temporary input and output file paths
  const tempDir = os.tmpdir();
  const inputPath = path.join(tempDir, `input-${Date.now()}`);
  const outputPath = path.join(tempDir, `output-${Date.now()}.mp3`);

  try {
    // Write the input file to disk
    const buffer = await audioFile.arrayBuffer();
    await fs.writeFile(inputPath, Buffer.from(buffer));

    // Convert using ffmpeg
    await new Promise<void>((resolve, reject) => {
      ffmpeg(inputPath)
        .output(outputPath)
        .on("end", () => resolve())
        .on("error", (err) => reject(err))
        .run();
    });

    // Read the converted file
    const outputBuffer = await fs.readFile(outputPath);

    // Create a new File object
    const convertedFile = new File([outputBuffer], `converted.mp3`, {
      type: "audio/mpeg",
    });

    return convertedFile;
  } finally {
    // Clean up temporary files
    try {
      await fs.unlink(inputPath).catch(() => {});
      await fs.unlink(outputPath).catch(() => {});
    } catch (error) {
      console.error("Error cleaning up temp files:", error);
    }
  }
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
    // Using the audio-enabled model to get both text and audio responses
    const completion = await client.chat.completions.create({
      model: "gpt-4o-mini-audio-preview",
      messages: openAiMessages,
      modalities: ["text", "audio"],
      audio: {
        format: "mp3",
        voice: "alloy",
      },
      max_completion_tokens: 2048,
    });

    // Create the AI message from the text response
    const aiMessage: Message = {
      type: MessageType.AiMessage,
      content:
        completion.choices[0].message.content ||
        completion.choices[0].message.audio?.transcript ||
        "",
      timestamp: new Date().toISOString(),
    };

    // Send the text response first
    onTextReady(aiMessage);
    console.log("Text response sent:", aiMessage.content);

    // If audio is available, send the base64 data directly
    if (completion.choices[0].message.audio?.data) {
      const base64Data = completion.choices[0].message.audio.data;
      onAudioReady(base64Data);
      console.log("Audio response sent");
    } else {
      console.log("No audio response");
    }

    return {
      aiMessage,
      audio: completion.choices[0].message.audio,
      usage: completion.usage ? analyzeUsage(completion.usage) : undefined,
    };
  } catch (error) {
    console.error("Error in audioChatStream:", error);
    throw error;
  }
};
