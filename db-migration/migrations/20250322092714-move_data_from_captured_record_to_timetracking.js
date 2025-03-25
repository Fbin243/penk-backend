const { ObjectId } = require('mongodb');

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
        const capturedRecords = await db.collection("captured_records").find().toArray();
        console.log(capturedRecords);
        const timeTrackings = [];
        capturedRecords.forEach((capturedRecord) => {
          const newTimetrackings = capturedRecord.time_trackings?.map((timetracking) => {
            const { time, ...rest } = timetracking;

            return {
              character_id: capturedRecord.metadata.character_id,
              _id: new ObjectId(),
              ...rest,
            };
          }) || [];

          timeTrackings.push(...newTimetrackings);
        });
        await db.collection("time_trackings").insertMany(timeTrackings);
        await db.collection("captured_records").drop();
        await db.collection("analytic_changelog").drop();
      });
    } finally {
      session.endSession();
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
