import ffmpeg from "fluent-ffmpeg";
import { promises as fs } from "fs";
import { parseBuffer } from "music-metadata";
import os from "os";
import path from "path";

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

    const outputBuffer = await fs.readFile(outputPath);

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

export const getAudioDuration = async (mp3Buffer: Buffer): Promise<number> => {
  const metadata = await parseBuffer(mp3Buffer, "audio/mpeg");
  return metadata.format.duration ?? 0;
};
