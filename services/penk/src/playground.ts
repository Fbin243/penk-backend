import "./bootstrap";

import chalk from "chalk";
import { createReadStream, writeFileSync } from "fs";

import { chat, transcribeAudio } from "./services/ai/utils";

const main = async () => {
  console.log(chalk.blue("Starting transcription..."));
  const transcriptionStartTime = Date.now();

  const audio = createReadStream("audio/PenK, bạn biết gì về tui.wav");
  const transcription = await transcribeAudio(audio);

  const transcriptionEndTime = Date.now();
  const transcriptionDuration = (transcriptionEndTime - transcriptionStartTime) / 1000;
  console.log(chalk.blue(`Transcription completed in ${transcriptionDuration.toFixed(2)} seconds`));

  console.log(chalk.green("Transcription:"));
  console.log(transcription);
  console.log();

  console.log(chalk.blue("Starting chat..."));
  const chatStartTime = Date.now();

  const response = await chat({
    userData: "Tui tên là Hiếu, tui rất đẹp trai",
    history: [],
    newMessage: transcription.text,
    voiceMode: true,
  });

  const chatEndTime = Date.now();
  const chatDuration = (chatEndTime - chatStartTime) / 1000;
  console.log(chalk.blue(`Chat completed in ${chatDuration.toFixed(2)} seconds`));

  if (response.audio) {
    writeFileSync("audio/response.wav", Buffer.from(response.audio.data, "base64"), {
      encoding: "utf-8",
    });
  }

  console.log(chalk.green("Response:"));
  console.log(response.aiMessage.content);
  console.log(chalk.green("Usage:"));
  console.log(response.usage);
  console.log();

  const totalCost = transcription.cost + (response.usage?.totalCost ?? 0);
  console.log(chalk.green("Total cost:"));
  console.log(totalCost);
  console.log();

  const totalDuration = transcriptionDuration + chatDuration;
  console.log(chalk.green("Total execution time:"));
  console.log(`${totalDuration.toFixed(2)} seconds`);
  console.log();
};

main();
