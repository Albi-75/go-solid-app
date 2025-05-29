import { createResource } from "solid-js";

type HelloResponse = {
  message: string;
};

const fetchHello = async (): Promise<HelloResponse> => {
  const res = await fetch("/api/hello");
  if (!res.ok) throw new Error("Failed to fetch");
  return res.json();
};

export default function HelloMessage() {
  const [data] = createResource<HelloResponse>(fetchHello);

  return (
    <div class="p-6 max-w-md mx-auto mt-10 bg-white rounded-2xl shadow-md">
      <h2 class="text-xl font-bold text-gray-800 mb-2">Message from Go API:</h2>
      {data.loading && <p class="text-gray-500">Loading...</p>}
      {data.error && <p class="text-red-500">Error: {data.error.message}</p>}
      {data() && <p class="text-green-600 font-mono">{data()?.message || "undefined"}</p>}
    </div>
  );
}
