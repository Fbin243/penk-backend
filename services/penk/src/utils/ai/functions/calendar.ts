import { zodFunction } from "openai/helpers/zod";
import { rrulestr } from "rrule";
import { z } from "zod";

import { getCalendarEvents } from "../../googleapis";
import { LinkedAccountType, Tool } from "../../types";
import { getLinkedAccounts } from "./../../database/utils";
import { SharedDescription } from "./shared";

const getCalendarEventsParams = z.object({
  profileId: z.string().describe(SharedDescription.profileId),
  timeMin: z
    .string()
    .describe(
      "Start time for the calendar query in ISO 8601 format (e.g., '2025-04-22T00:00:00Z'). Should be within 1 month of the current date.",
    ),
  timeMax: z
    .string()
    .describe(
      "End time for the calendar query in ISO 8601 format (e.g., '2025-04-29T23:59:59Z'). Should be within 1 month of the current date and after timeMin.",
    ),
  timezone: z.string().describe(SharedDescription.timezone),
  locale: z.string().describe(SharedDescription.locale),
});

export const functionGetCalendarEvents = async (props: {
  profileId: string;
  timeMin: string;
  timeMax: string;
  timezone: string;
  locale: string;
}) => {
  console.log(`[Tool: ${Tool.GetCalendarEvents}]`);
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

  return result;
};

export const toolGetCalendarEvents = zodFunction({
  name: Tool.GetCalendarEvents,
  description:
    "Retrieves the user's calendar events across all their linked calendars. Use this tool when the user asks about their schedule, upcoming meetings, events, or availability. The function defaults to showing events for the next 7 days if no specific timeframe is mentioned. Time queries are limited to a 1-month range in either the past or future. Results include event titles, times, locations, participants, and other relevant details from the user's calendars.",
  parameters: getCalendarEventsParams,
});
