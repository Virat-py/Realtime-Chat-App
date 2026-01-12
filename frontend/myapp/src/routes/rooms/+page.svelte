<script>
  export let data;
  console.log(data.rooms[0].room_name);
</script>

{#each data.rooms as room}
  <div>{room.room_id}) {room.room_name}</div>
{/each}

