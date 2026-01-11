<script>
    import { goto } from '$app/navigation';

    let name = "";
    let password = "";
    let token = "";

    async function register_user() {
        const res = await fetch("http://localhost:8080/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                UserID: name,
                Password: password,
            }),
        });

        if (res.ok){
          const data = await res.json()
          token = data.token;
          localStorage.setItem('token', token);
          redirect_user()
        }
    }

    async function login_user() {
        const res = await fetch("http://localhost:8080/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                UserID: name,
                Password: password,
            }),
        });
        console.log(res)
        
        if (res.ok){
          const data = await res.json()
          token = data.token;
          localStorage.setItem('token', token);
          redirect_user()
        }
    }

    function redirect_user(){
      goto("/rooms")
    }
</script>

<div>
    <input bind:value={name} placeholder="username" />
    <input bind:value={password} placeholder="password" />
</div>

<button on:click={register_user}>Register</button>
<button on:click={login_user}>Login</button>
