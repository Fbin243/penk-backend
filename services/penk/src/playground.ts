import "./bootstrap";

import { getPenKData } from "./utils/database/utils";

const testUserId = "6735a19cc0e37098e0286d6b";

const main = async () => {
  const userData = await getPenKData(testUserId);
  console.log("[PenK Context]");
  console.dir(userData, { depth: null, colors: true });
  console.log();
};

main();
