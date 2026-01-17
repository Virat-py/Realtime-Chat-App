<script>
    import { goto } from "$app/navigation";
    export let data;
    let newRoomName = "";

    async function create_new_room() {
        const res = await fetch("http://localhost:8080/create_room", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                room_name: newRoomName,
            }),
        });

        if (res.ok) {
            window.location.reload();
        }
        if (!res.ok) {
            const err = await res.json();
            alert(err.message || "Request failed");
            return;
        }
    }

    async function redirectToRoom(room_id) {
        try {
            goto(`room/${room_id}`);
        } catch {
            alert("Request failed");
            return;
        }
    }
</script>

<!-- <input bind:value={newRoomName} placeholder="Enter New Room Name" />
<button on:click={create_new_room}>Create New Room</button> -->
<div
    class="grid gap-6 p-6"
    style="grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));"
>
    <div class="relative flex items-center">
        <input
            bind:value={newRoomName}
            name="username"
            type="text"
            required
            class="w-full text-slate-900 text-sm border border-slate-300 px-4 py-3 pr-8 rounded-md outline-blue-600"
            placeholder="New Room Name"
        />
        
    </div>
    <button
        type="button"
        class="w-full px-4 py-1.5 rounded-md text-white text-sm font-medium tracking-wider border-none outline-none bg-blue-600 hover:bg-blue-700 cursor-pointer"
        on:click={create_new_room}
        >Add New Room</button
    >
</div>

<div
    class="grid gap-6 p-6"
    style="grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));"
>
    {#each data.rooms as room}
        <div
            class="bg-white border border-gray-200 shadow-md w-full max-w-sm rounded-lg overflow-hidden mx-auto mt-4"
        >
            <div class="p-4">
                <div>
                    <h3 class="text-lg font-semibold">{room.room_name}</h3>
                </div>
                <button
                    type="button"
                    class="mt-6 px-5 py-2 rounded-md text-white text-sm font-medium tracking-wider border-none outline-none bg-blue-600 hover:bg-blue-700 cursor-pointer"
                    on:click={() => redirectToRoom(room.room_id)}
                    >Chat Now</button
                >
            </div>
        </div>
    {/each}
</div>
