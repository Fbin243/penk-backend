const collections = require("../collections");

module.exports = {
  async up(db, client) {
    const session = client.startSession();
    try {
      // Rename custome_metrics -> categories, custom_metric_id -> category_id
      await session.withTransaction(async () => {
        const capturedRecords = await db.collection(collections.CapturedRecord).find({
          custom_metrics: { $exists: true },
        }).toArray();

        console.log("capturedRecords", capturedRecords);
        const newCapturedRecords = capturedRecords.map((record) => {
          const { custom_metrics, ...recordRest } = record;

          return {
            ...recordRest,
            categories: custom_metrics,
            time_trackings: record.time_trackings.map((timetracking) => {
              const { custom_metric_id, ...timetrackingRest } = timetracking;

              return {
                ...timetrackingRest,
                category_id: custom_metric_id,
              }
            }),
          }
        });

        console.log("newCapturedRecords", JSON.stringify(newCapturedRecords, null, 2));

        // Remove old records and insert new ones
        await db.collection(collections.CapturedRecord).deleteMany({ custom_metrics: { $exists: true } });
        await db.collection(collections.CapturedRecord).insertMany(newCapturedRecords);
      });
    } finally {
      await session.endSession();
    }
  },

  async down(db, client) { }
};
