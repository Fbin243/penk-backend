const { ObjectId } = require('mongodb');
const { mongodb } = require('../migrate-mongo-config');

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
        // Drop collections of template
        await db.collection('templates').drop();
        await db.collection('template_categories').drop();
        await db.collection('template_topics').drop();

        // Normalize category and metric from character
        const metricArr = []
        const categoryArr = []
        const filter = {}
        const characters = await db.collection('characters').find(filter).toArray();
        const newCharacters = characters.map((character) => {
          const { categories, metrics, total_focused_time, gender, tags, ...characterRest } = character;
          const newCategories = categories?.map((category) => {
            const { metrics, ...categoryRest } = category;
            const newMetrics = metrics?.map((metric) => {
              return {
                ...metric,
                category_id: category._id,
                character_id: character._id,
              };
            }) || [];
            metricArr.push(...newMetrics);
            return {
              ...categoryRest,
              character_id: character._id,
            };
          }) || [];
          categoryArr.push(...newCategories);

          const newMetrics = metrics?.map((metric) => {
            return {
              ...metric,
              character_id: character._id,
            };
          }) || [];
          metricArr.push(...newMetrics);

          return {
            ...characterRest,
          }
        });

        console.log('categoryArr: ', categoryArr);
        console.log('metricArr: ', metricArr);
        await db.collection('characters').deleteMany(filter);
        await db.collection('characters').insertMany(newCharacters);
        await db.collection('categories').insertMany(categoryArr);
        await db.collection('metrics').insertMany(metricArr);
        await db.collection('goals').drop();
        await db.collection('snapshots').drop();
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
    await db.collection('categories').drop();
    await db.collection('metrics').drop();
  }
};
