<script>
    import { onMount, onDestroy } from "svelte";
    export let data;
    
    console.log(data);
    let socket;
    let currMessage = "";

    // 👇 initialize from server data
    let messages = data.data ?? [];

    const roomID = data.room_id;

    onMount(() => {
        socket = new WebSocket(`ws://localhost:8080/ws?room_id=${roomID}`);

        socket.onmessage = (event) => {
            const msg = JSON.parse(event.data);

            // ✅ THIS triggers re-render
            messages = [...messages, msg];
        };

        socket.onerror = (err) => {
            console.error("WebSocket error", err);
        };
    });

    onDestroy(() => {
        socket?.close();
    });

    function sendMessage() {
        if (!currMessage.trim()) return;

        socket.send(
            JSON.stringify({
                Message: currMessage,
            }),
        );

        currMessage = "";
    }
</script>


<h3>Welcome to Room {roomID}</h3>


<div class="p-5">
    {#each messages as msg}
        <div class="p-1">
            <li class="max-w-lg ms-auto flex justify-end gap-x-2 sm:gap-x-4">
                <div class="grow text-end space-y-3">
                    <!-- Card -->
                    <div
                        class="inline-block bg-blue-600 rounded-2xl p-4 shadow-2xs"
                    >
                        <p class="text-sm text-white">
                            <b>{msg.user_id}:</b>
                            {msg.message}
                        </p>
                    </div>
                    <!-- End Card -->
                </div>
            </li>
        </div>
    {/each}
</div>
<div class="p-1">
    <li class="max-w-lg flex gap-x-2 sm:gap-x-4">
        <!-- Card -->
        <div class="bg-white border border-gray-200 rounded-2xl p-4 space-y-3">
            <p class="text-sm text-gray-800">
                hey
            </p>
            
        </div>
        <!-- End Card -->
    </li>
</div>

<input bind:value={currMessage} placeholder="Type..." />
<button on:click={sendMessage}>Send</button>
