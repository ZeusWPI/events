# Events

If you’re reading this, I assume you’re interested in how the code works and maybe even thinking about contributing :tada:.

This application handles quite a lot. Combined with the fact that Golang is a new language for many, the codebase can seem daunting at first.
Hopefully this readme will help you get your feet wet.

## Structure

Just like a good lasagna, this project is split into layers. Each with its own responsibility.

Let’s go through them, starting from the deepest layer: the database.

1. **Database** \
It all starts with database queries. The queries are written by hand and can be found in the [db/queries directory](../db/queries/). \
To add type safety, we generate Golang functions from these queries using [sqlc](https://sqlc.dev/). The generated functions accept the required arguments to execute the query.

2. **Repository** \
While sqlc is a great tool it, it doesn't alway generate ideal Golang structs. For this reason, we maintain our own [models directory](./db/model/), which contains internal representations of our database models. \
The repository layer acts as a bridge:

    - It converts between sqlc-generated structs and our models (both ways).
    - It ensures that the rest of the application only works with models.

3. **Service** \
Sometimes the API needs data in a slightly different format. That's why we use [DTOs](./server/dto/). \
The service layers handles:

    - Converting between DTOs and models.
    - Applying business logic.
    - Retuning results back to the API as DTOs.

**Example flow (creating an object from the frontend)**

1. The API receives the request, parses it into a DTO, and passes it to the service layer.
2. The service layer converts the DTO to a model and sends it to the repository.
3. The repository converts it to a sqlc struct, inserts it into the database, adds an ID, and returns it.
4. The service layer converts the new object back into a DTO and returns it to the API.
5. The API sends the DTO back to the frontend.

Is this the most efficient way? No. \
Does it create duplicate code accross layers? Yes.

But it enforces a **clear separation of responsibilities**, making debugging and feature development much easier.

## Tasks

A task is a function that runs in the background, either periodically or on demand. \
Tasks come by default with:

- Logs
- On-demand execution
- Status tracking

All tasks are managed by a central [Task Manager](./task/manager.go), which keeps track of their status, logs, etc.

To add a task:

1. Implement the `Task` interface.
2. Register your struct to the Task Manager.

The `Task` interface defines the required methods so the manager can expose the data and run the task. \
**Note:** the Task Manager itself doesn’t log execution. Instead, the task package provides a `NewTask` wrapper that adds logging for you.

**The Task Manager does not keep track of state!** \
This means that you need to reschedule tasks yourself.

- For periodic tasks this simply means registering your struct every startup
- For one time tasks you should keep track if the task needs to be rescheduled.

The logs are saved based on the given `TaskUID`.
Never change this after registering for the first time.
If you do change it then the previous logs will no longer be shown in the frontend (but they will still be in the database).

The `Task` interface has a method `Name()` which determines how your task will be shown in the frontend.
You can change this value as much as you like

## Checks

A check belongs to an event and can be interpreted as a TODO for that event.
Checks can either be

- Automatic: e.g. is an announcement written for that event
- Manual: created by an user in the frontend

Automatic checks will change their statuses automatically, while for manual checks the user has to make status changes in the frontend.

The internal logic is very similar to the tasks.
We again have a `Check Manager` where you need to register your checks by implementing the `Check` interface.

The package registering a check is responsible for letting the manager know when the status is updated.
However the check manager will automatically set statuses to `TODOLate` if the check has a deadline and hasn't been completed in time.

## Member vs Organizer vs Board

The most confusing database structure is the one for **users**. The backend and frontend use slightly different definitions.

#### Backend

- **Member:** Anyone who logs in is added to the member table. This table mainly connects users to their Zauth ID (so we can link board members).
- **Organizer:** Connects a board member with an event, indicating responsibility for organizing it.
- **Board:** A board member is linked to an academic year. This means one of three things:
  1. The user appears on the website as a board member for that year.
  2. The user has the events_admin role for that year.
  3. (Development only) The user is granted access for testing.

**Authorization rules:** \
You are authorized to access the application if:

1. It’s the development environment or you have on these these roles:
    - bestuur
    - events_admin
2. You are in the board table.

The second condition explains why we add users to the board table in development, and why there’s a field `is_organizer` to indicate event-organizing rights.

We could remove the board table and only keep the users but we'd lose two key features:

- Assigning events to users who haven't logged in yet.
- Automatically resetting access each academic year.

#### Frontend

The frontend doesn’t need this complexity.
It only knows about the type Organizer, which combines board and member information.

The frontend only needs to know:

1. Who can organize events in a given year.
2. Who the logged-in user is.
