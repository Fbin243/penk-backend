import * as admin from "firebase-admin";
import fs from "fs";

const firebaseAdminConfig = JSON.parse(
  fs.readFileSync(`../../${process.env.FIREBASE_ADMIN}`, "utf8"),
);

admin.initializeApp({
  credential: admin.credential.cert(firebaseAdminConfig),
});

export const decodeFirebaseJwt = async (token: string): Promise<admin.auth.DecodedIdToken> => {
  return await admin.auth().verifyIdToken(token);
};
