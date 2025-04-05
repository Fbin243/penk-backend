import chalk from "chalk";
import { ChatCompletionTool } from "openai/resources/index.mjs";
import { FunctionDefinition } from "openai/resources/shared.mjs";
import { rrulestr } from "rrule";

import { OAuthTokenModel } from "../../database/mongo";
import { decrypt } from "../../encrypt";
import { getCalendarEvents, refreshToken } from "../../googleapis";
import { LinkedAccountType } from "../../types";
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
      locale: {
        type: "string",
        description: "Locale of the user",
      },
    },
    required: ["profileId", "timeMin", "timeMax", "locale"],
    additionalProperties: false,
  },
  strict: true,
};

export const functionGetCalendarEvents = async (props: {
  profileId: string;
  timeMin: string;
  timeMax: string;
  locale: string;
}) => {
  console.log(chalk.cyan(`[Tool: ${FunctionName.GetCalendarEvents}]`));
  console.dir(props, { depth: null, colors: true });
  console.log();

  const tokens = await OAuthTokenModel.find({
    profile_id: props.profileId,
    type: LinkedAccountType.GoogleCalendar,
  });

  const eventFetchPromises = tokens.map(async (token) => {
    try {
      const accessToken = await refreshToken(decrypt(token.refresh_token));

      if (!accessToken) {
        console.warn(`Failed to refresh token for account: ${token.email}`);
        return [];
      }

      const events = await getCalendarEvents({
        accessToken,
        timeMin: new Date(props.timeMin),
        timeMax: new Date(props.timeMax),
      });

      return events || [];
    } catch (err) {
      console.error(`Error fetching events for ${token.email}:`, err);
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

  return rruleResolvedEvents.map((e) => ({
    title: e.summary,
    description: e.description,
    start: e.start ? new Date(e.start).toLocaleString(props.locale) : undefined,
    end: e.end ? new Date(e.end).toLocaleString(props.locale) : undefined,
  }));
};

export const toolGetCalendarEvents: ChatCompletionTool = {
  type: "function",
  function: getCalendarEventsDefinition,
};
