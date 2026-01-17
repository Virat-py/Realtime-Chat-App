<script>
    import { goto } from "$app/navigation";

    let name = $state("");

    let password = $state("");
    let token = "";

    async function register_user() {
        const res = await fetch("http://localhost:8080/register", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                UserID: name.value,
                Password: password.value,
            }),
        });

        if (res.ok) {
            redirect_user();
        }
    }

    async function login_user() {
        console.log("im in login");
        const res = await fetch("http://localhost:8080/login", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                UserID: name.value,
                Password: password.value,
            }),
        });
        console.log(res);

        if (res.ok) {
            redirect_user();
        }
    }

    function redirect_user() {
        goto("/rooms");
    }
</script>

<div class="bg-gray-50">
    <div
        class="min-h-screen flex flex-col items-center justify-center py-6 px-4"
    >
        <div class="max-w-[480px] w-full">
            <div
                class="p-6 sm:p-8 rounded-2xl bg-white border border-gray-200 shadow-sm"
            >
                <h1 class="text-slate-900 text-center text-3xl font-semibold">
                    Register or Login
                </h1>
                <form class="mt-12 space-y-6">
                    <div>
                        <label
                            class="text-slate-900 text-sm font-medium mb-2 block"
                            >Username</label
                        >
                        <div class="relative flex items-center">
                            <input
                                bind:value={name}
                                name="username"
                                type="text"
                                required
                                class="w-full text-slate-900 text-sm border border-slate-300 px-4 py-3 pr-8 rounded-md outline-blue-600"
                                placeholder="Enter username"
                            />
                        </div>
                    </div>
                    <div>
                        <label
                            class="text-slate-900 text-sm font-medium mb-2 block"
                            >Password</label
                        >
                        <div class="relative flex items-center">
                            <input
                                bind:value={password}
                                name="password"
                                type="password"
                                required
                                class="w-full text-slate-900 text-sm border border-slate-300 px-4 py-3 pr-8 rounded-md outline-blue-600"
                                placeholder="Enter password"
                            />
                        </div>
                    </div>

                    <div class="!mt-12 space-y-3">
                        <button
                            on:click={register_user}
                            type="button"
                            class="w-full py-2 px-4 text-[15px] font-medium tracking-wide rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none cursor-pointer"
                        >
                            Register
                        </button>

                        <button
                            on:click={login_user}
                            type="button"
                            class="w-full py-2 px-4 text-[15px] font-medium tracking-wide rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none cursor-pointer"
                        >
                            Login
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>
</div>
