<template>
    <div class="container-fluid">
        <div style="text-align: center;">
            <p class="lead">Découvrez...explorez...définissez...votre langue</p>
        </div>

        <div style="text-align: center">
            <img style="max-height: 150px;" alt="Brand" src="../assets/images/dvlf_logo_medium_no_beta_transparent.png" v-if="atHome">
            <div class="row" style="margin-top: 30px;">
                <form class="col-xs-12 col-sm-offset-3 col-md-offset-4 col-sm-6 col-md-4" @submit.prevent @keyup.enter="search()">
                    <div class="input-group">
                        <input type="text" class="form-control" :name="queryTerm" placeholder="Tapez un mot..." aria-describedby="search" v-model="queryTerm">
                        <span class="input-group-btn">
                            <button class="btn btn-default" type="submit">Rechercher</button>
                        </span>
                    </div>
                </form>
            </div>
        </div>
        <div style="text-align: center; margin-top: 50px;" v-if="atHome">
            <h3>Notre mot du jour :
                <router-link :to="`/mot/${wordOfTheDay}`">{{wordOfTheDay}}</router-link>
            </h3>
        </div>
    </div>
</template>
<script>
import { EventBus } from "../main.js"

export default {
    name: "Search",
    data: function() {
        return {
            atHome: true,
            queryTerm: "",
            wordOfTheDay: "tradition",
            aPropos: false
        }
    },
    watch: {
        $route(to, from) {
            if (to.path === "/") {
                this.atHome = true
            }
        }
    },
    created() {
        if (this.$route.path.slice(0, 4) == "/mot") {
            this.atHome = false
        }
        let vm = this
        EventBus.$on("OffHome", function() {
            vm.atHome = false
        })
    },
    methods: {
        search() {
            this.$router.push(`/mot/${this.queryTerm}`)
            this.atHome = false
        }
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.form-control {
    font-size: 100%;
}
</style>