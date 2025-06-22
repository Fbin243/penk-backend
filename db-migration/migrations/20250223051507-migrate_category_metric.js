const cols = {
  Profile: "profiles",
  Character: "characters",
}

module.exports = {
  async up(db, client) {
    const session = client.startSession();
    try {
      await session.withTransaction(async () => {
        // ----------- Profile
        await db.collection(cols.Profile).updateMany(
          {},
          {
            $unset: {
              limited_character_number: "",
            },
          }
        );

        // -------------- Character
        const characters = await db.collection(cols.Character).find({
          custom_metrics: { $exists: true },
        }).toArray();

        // custom_metrics -> categories, properties -> metrics 
        const newCharacters = characters.map((character) => {
          console.log(JSON.stringify(character, null, 2), "\n");
          let metrics = [];
          let categories = character.custom_metrics?.map(metric => {
            metrics = metrics.concat(metric.properties?.map(property => {
              const { type: propType, ...propRest } = property;
              return {
                ...propRest,
                category_id: metric._id,
                value: parseFloat(property.value) || 0,
              }
            }) || []);

            const { properties, limited_property_number, time, ...metricRest } = metric;
            return {
              ...metricRest,
            }
          });

          const { custom_metrics, limited_metric_number, total_focused_time, ...characterRest } = character;
          return {
            ...characterRest,
            categories,
            metrics,
          }
        });

        console.log(JSON.stringify(newCharacters, null, 2));
        // Remove old characters and insert new ones
        await db.collection(cols.Character).deleteMany({ custom_metrics: { $exists: true } });
        if (newCharacters.length > 0) {
          await db.collection(cols.Character).insertMany(newCharacters);
        }
      });
    } finally {
      await session.endSession();
    }
  },

  async down(db, client) { },
};
