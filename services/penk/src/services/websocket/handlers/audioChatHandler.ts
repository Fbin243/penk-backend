import { WebSocket } from "ws";

import { audioChat, base64ToUploadable, transcribeAudio } from "../../../utils/ai";
import { convertAudioFormatToMp3 } from "../../../utils/audio";
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
      const startTime = Date.now();

      const { data, type } = JSON.parse(message.toString()) as Ws_Message;

      if (type === Ws_MessageType.ConfigAudioFormat) {
        clientAudioFormat = data;
      }

      if (type === Ws_MessageType.UploadAudio) {
        const uploadStartTime = Date.now();

        const filename = `audio-${context.profileId}`;

        let uploadableAudio: File;
        try {
          // Convert to MP3 if needed for transcription (OpenAI doesn't support m4a)
          if (clientAudioFormat === "m4a") {
            const m4aFile = base64ToUploadable(data, `${filename}.m4a`, "audio/x-m4a");
            const mp3File = await convertAudioFormatToMp3(m4aFile);
            uploadableAudio = mp3File;
          } else {
            const mp4File = base64ToUploadable(data, `${filename}.mp4`, "audio/mp4");
            const mp3File = await convertAudioFormatToMp3(mp4File);
            uploadableAudio = mp3File;
          }
          console.log(`Audio format processing completed in ${Date.now() - uploadStartTime}ms`);
        } catch (error) {
          console.error("Error preparing audio for transcription:", error);
          sendErrorResponse(ws, "Error processing audio format");
          return;
        }

        try {
          const processingStartTime = Date.now();

          // Get PenK data and chat history
          const [transcriptionResult, penkData, penkMessages] = await Promise.all([
            transcribeAudio(uploadableAudio),
            getPenKData(context.userId),
            getPenKMessages(context.profileId),
          ]);
          console.log(`Transcribed message: ${transcriptionResult.text}`);
          console.log(`Transcription completed in ${Date.now() - processingStartTime}ms`);

          const transcriptResultMessage: Ws_Message = {
            type: Ws_MessageType.TranscriptResult,
            data: transcriptionResult.text,
            timestamp: new Date().toISOString(),
          };
          ws.send(JSON.stringify(transcriptResultMessage));

          // Create user message from transcription
          const userMessage: Message = {
            type: MessageType.UserMessage,
            content: transcriptionResult.text,
            timestamp: new Date().toISOString(),
          };

          // Process with AI and generate response
          const aiStartTime = Date.now();

          const { cost } = await audioChat(
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
              const textResponseTime = Date.now();
              console.log(
                `[${context.email}] Text response ready in ${textResponseTime - aiStartTime}ms`,
              );

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
              const audioResponseTime = Date.now();
              console.log(
                `[${context.email}] Audio response ready in ${audioResponseTime - aiStartTime}ms`,
              );

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

          const totalCost = transcriptionResult.cost + cost;
          console.log("Total cost:", totalCost);

          const totalTime = Date.now() - startTime;
          console.log(`Total audio chat processing completed in ${totalTime}ms`);
        } catch (error) {
          console.error("Error processing audio:", error);
          sendInfoResponse(ws, Ws_InfoType.TranscriptionFailed);
          sendErrorResponse(ws, "Error processing audio");

          const totalTime = Date.now() - startTime;
          console.log(`Audio chat processing failed after ${totalTime}ms`);
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
