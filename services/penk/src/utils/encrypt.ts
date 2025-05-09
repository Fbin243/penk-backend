import crypto from "crypto";

// Creating a proper 32-byte key from the provided secret
const secretKey = crypto
  .createHash("sha256")
  .update(process.env.ENCRYPTION_SECRET_KEY || "penk_secret")
  .digest();

// Generate a new IV for each encryption operation
export const encrypt = (data: string) => {
  const iv = crypto.randomBytes(16);
  const cipher = crypto.createCipheriv("aes-256-cbc", secretKey, iv);
  const encrypted = cipher.update(data, "utf8", "hex") + cipher.final("hex");
  // Return both the IV and encrypted data so we can decrypt later
  return iv.toString("hex") + ":" + encrypted;
};

export const decrypt = (encryptedData: string) => {
  // Extract the IV and the encrypted data
  const parts = encryptedData.split(":");
  const iv = Buffer.from(parts[0], "hex");
  const encryptedText = parts[1];
  const decipher = crypto.createDecipheriv("aes-256-cbc", secretKey, iv);
  return decipher.update(encryptedText, "hex", "utf8") + decipher.final("utf8");
};
