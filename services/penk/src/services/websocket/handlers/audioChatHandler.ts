import chalk from "chalk";
import { WebSocket } from "ws";

import {
  audioChat,
  base64ToUploadable,
  convertAudioFormatToMp3,
  transcribeAudio,
} from "../../../utils/ai";
import { PenKMessageModel } from "../../../utils/database/mongo";
import { getPenKData, getPenKMessages } from "../../../utils/database/utils";
import {
  Message,
  MessageType,
  Ws_InfoType,
  Ws_Message,
  Ws_MessageType,
} from "../../../utils/types";
import { WebSocketContext } from "../types";
import { sendErrorResponse, sendInfoResponse } from "../utils";

export const handleAudioChat = (ws: WebSocket, context: WebSocketContext) => {
  console.log(`Audio chat connection established for user: ${context.email}`);

  let clientAudioFormat = "mp3";

  ws.on("message", async (message: string) => {
    try {
      const { data, type } = JSON.parse(message.toString()) as Ws_Message;

      if (type === Ws_MessageType.ConfigAudioFormat) {
        clientAudioFormat = data;
        console.log(`Audio format configured: ${clientAudioFormat}`);
      }

      if (type === Ws_MessageType.UploadAudio) {
        const filename = `audio-${context.profileId}`;

        let uploadableAudio;
        try {
          // Convert to MP3 if needed for transcription (OpenAI doesn't support m4a)
          if (clientAudioFormat === "m4a") {
            console.log("Converting input audio from m4a to mp3 for transcription");
            const m4aFile = base64ToUploadable(data, `${filename}.m4a`, "audio/x-m4a");
            const mp3File = await convertAudioFormatToMp3(m4aFile);
            uploadableAudio = mp3File;
          } else {
            uploadableAudio = base64ToUploadable(data, `${filename}.mp3`, "audio/mpeg");
          }
        } catch (error) {
          console.error("Error preparing audio for transcription:", error);
          sendErrorResponse(ws, "Error processing audio format");
          return;
        }

        try {
          // Transcribe the audio
          const transcriptionResult = await transcribeAudio(uploadableAudio);
          console.log(`Transcribed message: ${transcriptionResult.text}`);

          // Get PenK data and chat history
          const [penkData, penkMessages] = await Promise.all([
            getPenKData(context.userId),
            getPenKMessages(context.profileId),
          ]);

          // Create user message from transcription
          const userMessage: Message = {
            type: MessageType.UserMessage,
            content: transcriptionResult.text,
            timestamp: new Date().toISOString(),
          };

          // Process with AI and generate response
          const { usage } = await audioChat(
            {
              userData: JSON.stringify(penkData),
              history: penkMessages.map((message) => ({
                type: message.type,
                content: message.content,
                timestamp: message.timestamp.toISOString(),
              })),
              transcribedMessage: transcriptionResult.text,
            },
            // This callback is called when text response is ready
            (textResponse) => {
              // Save messages to database
              PenKMessageModel.insertMany([
                { ...userMessage, profile_id: context.profileId },
                { ...textResponse, profile_id: context.profileId },
              ]);

              // Send the text response to client
              const textResponseMessage: Ws_Message = {
                type: Ws_MessageType.TextChat,
                data: textResponse.content,
                timestamp: textResponse.timestamp,
              };
              ws.send(JSON.stringify(textResponseMessage));
            },
            // This callback is called when audio response is ready (always in mp3 format from OpenAI)
            async (base64Audio) => {
              try {
                // Send the audio response to client
                const audioResponseMessage: Ws_Message = {
                  type: Ws_MessageType.DownloadAudio,
                  data: base64Audio,
                  timestamp: new Date().toISOString(),
                };
                ws.send(JSON.stringify(audioResponseMessage));

                // Send completion notification
                sendInfoResponse(ws, Ws_InfoType.AudioStreamCompleted);
              } catch (error) {
                console.error("Error converting response audio format:", error);
                sendErrorResponse(ws, "Error processing audio response");
              }
            },
          );

          console.log(
            chalk.green("Total cost:", transcriptionResult.cost + (usage?.totalCost ?? 0)),
          );
        } catch (error) {
          console.error("Error processing audio:", error);
          sendInfoResponse(ws, Ws_InfoType.TranscriptionFailed);
          sendErrorResponse(ws, "Error processing audio");
        }
      }
    } catch (error) {
      console.error("Error with audio message:", error);
      sendErrorResponse(ws, "Error with audio message");
    }
  });

  ws.on("close", () => {
    console.log("Audio chat connection closed");
  });
};
