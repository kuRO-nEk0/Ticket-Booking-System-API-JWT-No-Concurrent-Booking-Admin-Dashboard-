# Ticket Booking Conversational Agent

A conversational agent for event discovery and booking, where finding and recommending an event is **one capability among several** rather than a rigid search flow.

```
you > looking for a jazz show this weekend in Guwahati

  THOUGHT     intent=event_search  branch=recommend  (city plus at least one constraint)
  ACTION      search_events({'city': 'Guwahati', 'tags': ['jazz'], 'date_range': 'this_weekend'})
  OBSERVATION 2 candidates

bot > Here are 2 options for this weekend:
        Midnight Jazz Quartet — The Blue Room
          acoustic • seated
```

---

## Problem

Recommend event tickets and guide users through booking from an open conversation rather than a static form. 
The user is talking to the assistant about their weekend plans; recommendations have to emerge naturally, and the actual booking must be handled securely without the LLM being in the write path.

## Approach

**The LLM does query understanding and explanation. The database does retrieval and booking.**

Passing hundreds of events into a prompt is slow, non-deterministic, and invents events that do not exist. So the model translates the request in and narrates the result out; a deterministic filter decides the availability.

The harder problem is **knowing when to recommend at all**. The system can fail three ways:

- **Too eager** — user asks about the weather, gets five outdoor events. Reads as advertising, and people stop trusting it.
- **Too passive** — user says "I want to see some standup tonight" and the agent just chats pleasantly. The entire point, missed.
- **Too early** — recommends knowing only the city. Returns a generic list no better than browsing, and burns the first impression.

Hence a router that classifies every turn into one of four branches, and only one of them touches the database.

## Event data design — three layers

The split exists because seat availability is a **dealbreaker** and "acoustic" is a **preference**.

| Layer | Contents | Mechanism |
|---|---|---|
| 1 — Hard constraints | city, category, dates, capacity, available seats | SQL `WHERE` / exact filter |
| 2 — Semantic profile | tags and mood classification (e.g. `chill`, `big_night`) | exact match / overlaps |
| 3 — Derived signals | near sold-out status | ranking/urgency flags |

**Invariant:** The agent uses a strict controlled vocabulary for tags and categories (defined in `core/agent/vocab.go`). 
If the user asks for a "relaxed evening", it maps to the `chill` mood, which safely translates to `acoustic`, `jazz`, or `seated` tags.

## System design

Each turn: **route → decide → act → compose.**

| Branch | When | Touches DB |
|---|---|---|
| `out_of_scope` | not an event/entertainment question | no |
| `answer_directly` | chitchat, or planning with no destination settled | no |
| `ask_one` | ambiguous request (e.g. "the jazz one" when there are two) | no |
| `recommend` | confident city/date **plus** at least one discriminator | yes |

Two design decisions worth naming:

**The LLM is never in the write path.** The agent can draft a `BookingDraft`, but it cannot create a finalized `Booking`. The actual booking is finalized by the canonical backend route relying on Optimistic Concurrency Control (`Seat.Version`). A second write path is a second chance to double-book.

**Inferred slots carry confidence.** Only confidently extracted slots reach the filter. Lower-confidence ones shape the reply so the user can correct them. 

**The relaxation loop** is what makes this an agent rather than a pipeline. An empty result set is an Observation: the agent gives up one constraint, retries, and says what it gave up — *"nothing tonight for standup, so I checked tomorrow."*

## Running it

```bash
# 1. Start the backend
cd backend
go mod tidy
go run ./cmd/api/main.go
```

```bash
# 2. Seed the database with spread, consistent events
cd backend
go run ./cmd/seed/main.go
```
*Note: The seed script is idempotent and safely guarded with a marker event. It generates 36 procedurally varied events including ambiguity fixtures.*

```bash
# 3. Start the frontend
cd frontend
npm install
npm run dev
```

## Structure

```
ticket-booking-system/
├── frontend/                The SvelteKit Web UI
│   └── src/                 Component library, dynamic routing, and animations
├── backend/                 The Go Backend
│   ├── cmd/
│   │   ├── api/main.go      The Fiber API entrypoint (Local)
│   │   └── seed/main.go     Database generation and seeding script
│   ├── api/index.go         Vercel Serverless entrypoint
│   └── core/
│       ├── models/          Data schemas (Events, Bookings, Agent Models)
│       ├── controllers/     API Logic and Seat Generation
│       ├── agent/           Controlled vocabulary and routing structures
│       └── repository/      PostgreSQL / SQLite Database logic
```

## Status and limitations

- **No payment gateway.** Read-only event browsing and dummy booking flow.
- **Seat selection is simplified.** The agent drafts a quantity, and the system assigns available seats rather than a visual seat-map selection.

