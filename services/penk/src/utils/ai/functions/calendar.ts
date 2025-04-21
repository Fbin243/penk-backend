import chalk from "chalk";
import { ChatCompletionTool } from "openai/resources/index.mjs";
import { FunctionDefinition } from "openai/resources/shared.mjs";
import { rrulestr } from "rrule";

import { getCalendarEvents } from "../../googleapis";
import { LinkedAccountType } from "../../types";
import { getLinkedAccounts } from "./../../database/utils";
import { FunctionName } from "./types";

const getCalendarEventsDefinition: FunctionDefinition = {
  name: FunctionName.GetCalendarEvents,
  description:
    "Fetch calendar events from multiple linked calendars. Defaults to the next 7 days. Limit queries to a 1-month range in the past or future.",
  parameters: {
    type: "object",
    properties: {
      profileId: {
        type: "string",
        description: "User's profile ID",
      },
      timeMin: {
        type: "string",
        description: "Start time (ISO 8601 format)",
      },
      timeMax: {
        type: "string",
        description: "End time (ISO 8601 format)",
      },
      timezone: {
        type: "string",
        description: "Timezone of the user",
      },
      locale: {
        type: "string",
        description: "Locale of the user",
      },
    },
    required: ["profileId", "timeMin", "timeMax", "timezone", "locale"],
    additionalProperties: false,
  },
  strict: true,
};

export const functionGetCalendarEvents = async (props: {
  profileId: string;
  timeMin: string;
  timeMax: string;
  timezone: string;
  locale: string;
}) => {
  console.log(chalk.cyan(`[Tool: ${FunctionName.GetCalendarEvents}]`));
  console.dir(props, { depth: null, colors: true });
  console.log();

  const linkedAccounts = (await getLinkedAccounts(props.profileId)).filter(
    (linkedAccount) => linkedAccount.type === LinkedAccountType.GoogleCalendar,
  );

  if (linkedAccounts.length === 0) {
    throw new Error(
      'No linked Google Calendar accounts found. Please ask user to go to App Settings, click on the "Linked Accounts" button, and link their Google Calendar account.',
    );
  }

  const eventFetchPromises = linkedAccounts.map(async (linkedAccount) => {
    try {
      const events = await getCalendarEvents({
        accessToken: linkedAccount.accessToken,
        timeMin: new Date(props.timeMin),
        timeMax: new Date(props.timeMax),
      });

      return events || [];
    } catch (err) {
      console.error(`Error fetching events for ${linkedAccount.email}:`, err);
      return [];
    }
  });

  const allEvents = (await Promise.all(eventFetchPromises)).flat().map((e) => ({
    id: e.id,
    summary: e.summary,
    description: e.description,
    start: e.start?.dateTime,
    end: e.end?.dateTime,
    recurrence: e.recurrence,
  }));

  const rruleResolvedEvents = allEvents.flatMap((event) => {
    if (event.recurrence && event.recurrence.length > 0) {
      try {
        const ruleString = event.recurrence[0].replace("RRULE:", "");
        const rule = rrulestr(ruleString, { dtstart: new Date(event.start!) });

        const occurrences = rule.between(new Date(props.timeMin), new Date(props.timeMax));

        return occurrences.map((occurrence) => ({
          ...event,
          id: `${event.id}-${occurrence.toISOString()}`,
          start: occurrence.toISOString(),
          end: new Date(
            occurrence.getTime() +
              (new Date(event.end!).getTime() - new Date(event.start!).getTime()),
          ).toISOString(),
        }));
      } catch (error) {
        console.error(`Failed to parse RRULE for event ${event.id}:`, error);
      }
    }

    return [event]; // Return one-time events as-is
  });

  const result = rruleResolvedEvents.map((e) => ({
    title: e.summary,
    description: e.description,
    start: e.start
      ? new Date(e.start).toLocaleString(props.locale, {
          timeZone: props.timezone,
        })
      : undefined,
    end: e.end
      ? new Date(e.end).toLocaleString(props.locale, {
          timeZone: props.timezone,
        })
      : undefined,
  }));

  // console.log(chalk.cyan(`[Tool: ${FunctionName.GetCalendarEvents}]`));
  // console.dir(result, { depth: null, colors: true });
  // console.log();

  return result;
};

export const toolGetCalendarEvents: ChatCompletionTool = {
  type: "function",
  function: getCalendarEventsDefinition,
};
