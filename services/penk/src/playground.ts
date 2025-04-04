import "./bootstrap";

import chalk from "chalk";

import { getPenKData } from "./utils/database/utils";

const main = async () => {
  const userData = await getPenKData("67ef91fbff399a433c6ddc88");
  console.log(chalk.green("[PenK Context]"));
  console.dir(userData, { depth: null, colors: true });
};

main();
