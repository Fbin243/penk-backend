export const SharedDescription = {
  profileId: "User's unique profile identifier needed to access their data",
  firebaseUID: "User's unique Firebase identifier needed to access their data",
  assignedCategoryId:
    "A category ID to assign the item to. Must be either a valid string ID or explicitly null. Empty strings ('') are not acceptable and should be converted to null before assignment.",
  timezone:
    "User's timezone in IANA format (e.g., 'America/New_York', 'Europe/London') to correctly display event times",
  locale: "User's locale preference (e.g., 'en-US', 'fr-FR') for proper date and time formatting",
  datetime:
    "Data in ISO format (e.g., '2025-04-22T00:00:00Z'). Always setup date time on UTC because the timezone is handled on the client side.",
  eisenHowerMatrix: `Task priority using Eisenhower Matrix:
    - 1: Urgent and important
    - 2: Important but not urgent
    - 3: Urgent but not important
    - 4: Neither urgent nor important
    `,
};
