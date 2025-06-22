const mongo = require("mongodb")
const cols = {
  Profile: "profiles",
  Character: "characters",
  Goal: "goals",
  Template: "templates",
  TemplateCategory: "template_categories",
  TemplateTopic: "template_topics",
  Snapshot: "snapshots",
}

module.exports = {
  /**
   * @param db {import('mongodb').Db}
   * @param client {import('mongodb').MongoClient}
   * @returns {Promise<void>}
   */
  async up(db, client) {
    const session = client.startSession();
    try {
      await session.withTransaction(async () => {
        // ----------- Profile
        await db.collection(cols.Profile).updateMany(
          {},
          {
            $unset: {
              available_snapshots: "",
              auto_snapshot: "",
            },
          }
        );

        // -------------- Character
        const filter = {};
        const characters = await db.collection(cols.Character).find(
          filter,
        ).toArray();

        // move the metrics to categories
        const newCharacters = characters.map((character) => {
          console.log(JSON.stringify(character, null, 2), "\n");
          let { vision, metrics, categories, ...characterRest } = character;

          const cateMap = {};
          if (categories) {
            for (const cate of categories) {
              cateMap[cate._id] = cate;
            }

            metrics = metrics.filter(metric => {
              if (cateMap[metric.category_id]) {
                if (!cateMap[metric.category_id].metrics) {
                  cateMap[metric.category_id].metrics = [];
                }

                cateMap[metric.category_id].metrics.push(metric);
                delete metric.category_id;
                return false;
              }

              delete metric.category_id;
              return true;
            }) || [];
          }

          console.log("character_rest: ", JSON.stringify(characterRest, null, 2), "\n");
          console.log("cate_map: ", JSON.stringify(Object.values(cateMap), null, 2), "\n");
          console.log("metrics: ", JSON.stringify(metrics, null, 2), "\n");

          return {
            ...characterRest,
            categories: Object.values(cateMap),
            metrics,
          }
        });

        console.log(JSON.stringify(newCharacters, null, 2));
        // Remove old characters and insert new ones
        await db.collection(cols.Character).deleteMany(filter);
        await db.collection(cols.Character).insertMany(newCharacters);

        // -------------- Goal
        await db.collection(cols.Goal).drop();

        // -------------- Snapshot
        await db.collection(cols.Snapshot).drop();

        // -------------- Templates
        await db.collection(cols.TemplateCategory).drop();
      });
    } finally {
      await session.endSession();
    }
  },

  /**
   * @param db {import('mongodb').Db}
   * @param client {import('mongodb').MongoClient}
   * @returns {Promise<void>}
   */
  async down(db, client) {
    // TODO write the statements to rollback your migration (if possible)
    // Example:
    // await db.collection('albums').updateOne({artist: 'The Beatles'}, {$set: {blacklisted: false}});
  }
};
