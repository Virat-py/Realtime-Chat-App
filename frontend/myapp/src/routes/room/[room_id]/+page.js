export const load = async ({ fetch, params }) => {
  const room_id = params.room_id;
  
  let data;
  const res = await fetch(`http://localhost:8080/room/${room_id}`, {
    method: "GET",
    credentials: "include",
  });
  if (res.ok) {
    data = await res.json();
  }
  else {
    return {};
  }
  return {data,room_id};

};
