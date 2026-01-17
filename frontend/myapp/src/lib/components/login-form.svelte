<!-- <script>
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Card from "$lib/components/ui/card/index.js";
    import { Input } from "$lib/components/ui/input/index.js";
    import {
        FieldGroup,
        Field,
        FieldLabel,
        FieldDescription,
    } from "$lib/components/ui/field/index.js";

    let {
        name = $bindable(),
        password = $bindable(),
        register,
        login,
    } = $props();

    function submitRegister() {
        console.log("child register click");
        register?.();
    }

    function submitLogin() {
        console.log("child login click");
        login?.();
    }
</script>

<Card.Root class="mx-auto w-full max-w-sm">
    <Card.Header>
        <Card.Title class="text-2xl">Register or Login</Card.Title>
    </Card.Header>
    <Card.Content>
        <FieldGroup>
            <Field>
                <FieldLabel>Username</FieldLabel>
                <Input bind:value={name} type="text" required />
            </Field>
            <Field>
                <div class="flex items-center">
                    <FieldLabel>Password</FieldLabel>
                </div>
                <Input bind:value={password} type="password" required />
            </Field>
            <Field>
                <Button type="button" class="w-full" on:click={submitRegister}
                    >Register</Button
                >
                <Button
                    type="button"
                    variant="outline"
                    class="w-full"
                    on:click={submitLogin}>Login</Button
                >
            </Field>
        </FieldGroup>
    </Card.Content>
</Card.Root> -->
