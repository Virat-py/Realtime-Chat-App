export const load = async ({ fetch }) => {
  const res = await fetch("http://localhost:8080/get_rooms", {
    method: "GET",
    credentials: "include",
  });
  if (res.ok) {
    const rooms = await res.json();
    return {
      rooms
    };
  }
  return {};
};
