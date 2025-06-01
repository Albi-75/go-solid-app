import { createSignal, createResource, For } from "solid-js";

type Todo = {
  id: number;
  description: string;
  completed: boolean;
};

const fetchTodos = async (): Promise<Todo[]> => {
  const res = await fetch("/api/todos");
  if (!res.ok) throw new Error("Failed to fetch todos");
  return res.json();
};

export default function App() {
  const [text, setText] = createSignal("");
  const [todos, { refetch }] = createResource(fetchTodos);

  const deleteTodo = async (id: number) => {
    await fetch(`/api/todos/${id}`, { method: "DELETE" });
    refetch();
};

  const addTodo = async (e: Event) => {
    e.preventDefault();
    if (!text().trim()) return;
    await fetch("/api/todos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ description: text() }),
    });
    setText("");
    refetch();
  };

  const markDone = async (id: number) => {
    await fetch(`/api/todos/${id}/done`, { method: "PUT" });
    refetch();
  };

  return (
    <div class="min-h-screen bg-gray-100 flex flex-col items-center justify-center">
      <div class="bg-white rounded shadow-lg p-8 w-full max-w-md">
        <h1 class="text-2xl font-bold mb-6 text-center text-blue-600">My Todos</h1>
        <form onSubmit={addTodo} class="flex mb-6 gap-2">
          <input
            class="flex-1 rounded border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
            value={text()}
            onInput={e => setText((e.target as HTMLInputElement).value)}
            placeholder="What needs to be done?"
          />
          <button
            type="submit"
            class="bg-blue-600 text-white rounded px-4 py-2 hover:bg-blue-700 transition"
          >
            Add
          </button>
        </form>
        <ul class="space-y-3">
          <For each={todos()}>
            {todo => (
              <li class="flex items-center justify-between bg-gray-50 border border-gray-200 rounded px-4 py-2">
                <span
                  class={`flex-1 ${todo.completed ? "line-through text-gray-400" : ""}`}
                >
                  {todo.description}
                </span>
                {!todo.completed && (
                  <button
                    class="ml-4 text-sm bg-green-500 text-white px-3 py-1 rounded hover:bg-green-600 transition"
                    onClick={() => markDone(todo.id)}
                  >
                    Mark completed
                  </button>
                )}
                {todo.completed && (
                  <span class="ml-4 text-green-600 font-semibold text-xs">Done</span>
                )}
                <button
                  class="ml-2 text-sm bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600 transition"
                  onClick={() => deleteTodo(todo.id)}
                >
                  Delete
                </button>
              </li>
            )}
          </For>
        </ul>
      </div>
    </div>
  );
}