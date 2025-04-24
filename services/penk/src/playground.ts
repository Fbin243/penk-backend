import "./bootstrap";

import { HabitModel } from "./utils/database/mongo";
import { getPenKData } from "./utils/database/utils";

const testUserId = "6735a19cc0e37098e0286d6b";
const testProfileId = "6735a243c0e37098e0286d6c";

const main = async () => {
  const userData = await getPenKData(testUserId);
  console.log("[PenK Context]");
  console.dir(userData, { depth: null, colors: true });
  console.log();

  const query: Record<string, unknown> = { character_id: testProfileId };
  // if (props.categoryId) {
  //   query.category_id = props.categoryId;
  // }
  // if (props.name) {
  //   query.name = { $regex: props.name, $options: "i" };
  // }

  const habits = await HabitModel.find(query);

  console.log("[Habits]");
  console.dir(habits, { depth: null, colors: true });
  console.log();
};

main();
