<template>
    <div class='container'>
        <form>
            <div>
                <span class="grid field">
                    <label class="col-fixed" style="width:100px" for="username">Email</label>
                    <InputText id="username" type="text" v-model="authorizeForm.email"></InputText>
                </span>
                <span class="grid field">
                    <label class="col-fixed" style="width:100px" for="password">Password</label>
                    <InputText id="password" type="password" v-model="authorizeForm.password"></InputText>
                </span>
                <div class="formgroup-inline">
                    <span class="field">
                        <Button type="button" class="p-button-success" @click="onAuthorize">Login</Button>
                    </span>
                    <span class="field">
                        <Button type="button" @click="onRegistrate">Registrate</Button>
                    </span>
                </div>
            </div>
        </form>
    </div>
    <Dialog header="Registration" v-model:visible="showRegistrate" modal="true" @hide="initForm">
        <div class="formgroup">
            <span class="grid field">
                <label class="col-fixed" style="width:100px" for="name">Name</label>
                <InputText id="name" type="text" v-model="registrateForm.name" autofocus></InputText>
            </span>
            <span class="grid field">
                <label class="col-fixed" style="width:100px" for="surname">Surname</label>
                <InputText id="surname" type="text" v-model="registrateForm.surname"></InputText>
            </span>
            <span class="grid field">
                <label class="col-fixed" style="width:100px" for="age">Age</label>
                <InputNumber id="age" type="text" v-model="registrateForm.age"/>
            </span>
            <span class="grid field">
                <label class="col-fixed" style="width:100px" for="sex">Sex</label>
                <InputText id="sex" type="text" v-model="registrateForm.sex"></InputText>
            </span>
            <span class="grid field">
                <label class="col-fixed" style="width:100px" for="interests">Interests</label>
                <InputText id="interests" type="text" v-model="registrateForm.interests"></InputText>
            </span>
            <span class="grid field">
                <label class="col-fixed" style="width:100px" for="city">City</label>
                <InputText id="city" type="text" v-model="registrateForm.city"></InputText>
            </span>
            <span class="grid field">
                <label class="col-fixed" style="width:100px" for="email">Email</label>
                <InputText id="email" type="email" v-model="registrateForm.email"></InputText>
            </span>
            <span class="grid field">
                <label class="col-fixed" style="width:100px" for="password">Password</label>
                <InputText id="password" type="password" v-model="registrateForm.password"></InputText>
            </span>
        </div>
        <template #footer>
            <Button label="Cancel" icon="pi pi-times" class="p-button-text" @click="onClose"/>
            <Button label="Registrate" icon="pi pi-check" @click="onSubmit"/>
        </template>
    </Dialog>
</template>

<script>
import axios from 'axios';

export default {
    name: 'IndexPage',
    data() {
        return {
            registrateForm: {
                name: null,
                surname: null,
                age: null,
                sex: null,
                interests: null,
                city: null,
                email: null,
                password: null,
            },
            authorizeForm: {
                email: null,
                password: null,
            },
            showRegistrate: false,
        };
    },
    methods: {
        authorize(payload) {
            const path = `${process.env.VUE_APP_API_HOST || ''}/api/authorize`;
            // eslint-disable-next-line
            console.log(path);
            axios.post(path, payload)
                .then((result) => {
                    localStorage.setItem('token', result.data.token);
                    this.$router.push('/home');
                })
                .catch(() => {
                },
                );
        },
        registrate(payload) {
            const path = `${process.env.VUE_APP_API_HOST || ''}/api/registrate`;
            axios.post(path, payload)
                .catch((error) => {
                    // eslint-disable-next-line
                        console.log(error);
                    },
                );
        },
        initForm() {
            this.registrateForm.name = null;
            this.registrateForm.surname = null;
            this.registrateForm.age = null;
            this.registrateForm.sex = null;
            this.registrateForm.interests = null;
            this.registrateForm.city = null;
            this.registrateForm.email = null;
            this.registrateForm.password = null;
        },
        onAuthorize(evt) {
            evt.preventDefault();
            const payload = this.authorizeForm;
            this.authorize(payload);
        },
        onSubmit(evt) {
            evt.preventDefault();
            const payload = this.registrateForm;
            this.registrate(payload);
            this.showRegistrate = false;
        },
        onClose(evt) {
            evt.preventDefault();
            this.showRegistrate = false;
        },
        onRegistrate(evt) {
            evt.preventDefault();
            this.showRegistrate = true;
        },
    },
};
</script>

<style>

</style>
