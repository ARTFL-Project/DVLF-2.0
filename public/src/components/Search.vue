<template>
    <div class="container-fluid">
        <div style="text-align: center; margin-top: 30px;">
            <p class="lead">Découvrez...explorez...définissez...votre langue</p>
        </div>

        <div style="text-align: center">
            <img style="max-height: 150px;" alt="Brand" src="../assets/images/dvlf_logo_medium_no_beta_transparent.png" v-if="atHome">
            <b-row>
                <b-col sm="8" offset-sm="2" md="6" offset-md="3" lg="4" offset-lg="4" align-self="center">
                    <form @submit.prevent @keyup.enter="search()">
                        <b-input-group class="shadow-sm">
                            <b-form-input autocomplete="off" :name="queryTerm" placeholder="Tapez un mot..." aria-describedby="search" v-model="queryTerm" @input="onChange" @keyup.down.native="onArrowDown" @keyup.up.native="onArrowUp" @keyup.enter.native="onEnter"></b-form-input>
                            <b-input-group-append>
                                <b-btn variant="primary">Rechercher</b-btn>
                            </b-input-group-append>
                        </b-input-group>
                        <ul id="autocomplete-results" v-show="isOpen" class="autocomplete-results">
                            <li class="loading" v-if="isLoading">
                                Loading results...
                            </li>
                            <li v-else v-for="(result, i) in autoCompleteResults" :key="i" @click="setResult(result.headword)" class="autocomplete-result" :class="{ 'is-active': i === arrowCounter }" v-html="result.headword">
                            </li>
                        </ul>
                    </form>
                </b-col>
            </b-row>
        </div>
        <div style="text-align: center; margin-top: 10px;" v-if="atHome">
            <h4>Notre mot du jour :
                <router-link :to="`/mot/${wordOfTheDay}`">{{wordOfTheDay}}</router-link>
            </h4>
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
            typed: "",
            wordOfTheDay: "tradition",
            aPropos: false,
            isOpen: false,
            autoCompleteResults: ["amour", "gloire", "beauté"],
            isLoading: false,
            arrowCounter: 0
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
        this.$http(`${this.$globalConfig.apiServer}/api/wordoftheday`).then(
            function(response) {
                vm.wordOfTheDay = response.data
            }
        )
    },
    methods: {
        search() {
            this.$router.push(`/mot/${this.queryTerm}`)
            this.atHome = false
            this.isOpen = false
        },
        onChange() {
            // Let's warn the parent that a change was made
            this.isLoading = true
            this.typed = this.queryTerm
            if (this.typed.length > 1) {
                var vm = this
                this.$http
                    .get(
                        `${this.$globalConfig.apiServer}/api/autocomplete/${
                            this.typed
                        }`
                    )
                    .then(function(response) {
                        vm.autoCompleteResults = response.data.headwords
                        vm.isOpen = true
                        vm.isLoading = false
                    })
            }
        },

        filterResults() {
            // first uncapitalize all the things
            this.autoCompleteResults = this.items.filter(item => {
                return (
                    item.toLowerCase().indexOf(this.queryTerm.toLowerCase()) >
                    -1
                )
            })
        },
        setResult(result) {
            this.queryTerm = result.replace(/<[^>]+>/g, "")
            this.isOpen = false
        },
        onArrowDown(evt) {
            if (this.arrowCounter < this.autoCompleteResults.length) {
                this.arrowCounter = this.arrowCounter + 1
            }
        },
        onArrowUp() {
            if (this.arrowCounter > 0) {
                this.arrowCounter = this.arrowCounter - 1
            }
        },
        onEnter() {
            this.queryTerm = this.autoCompleteResults[
                this.arrowCounter
            ].headword.replace(/<[^>]+>/g, "")
            this.arrowCounter = -1
        },
        handleClickOutside(evt) {
            if (!this.$el.contains(evt.target)) {
                this.isOpen = false
                this.arrowCounter = -1
            }
        }
    },
    mounted() {
        document.addEventListener("click", this.handleClickOutside)
    },
    destroyed() {
        document.removeEventListener("click", this.handleClickOutside)
    }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.container-fluid {
    font-size: 90%;
}
.lead {
    margin-top: 15px;
}
.autocomplete {
    position: relative;
}

.autocomplete-results {
    padding: 0;
    margin: 0;
    border: 1px solid #eeeeee;
    max-height: 200px;
    overflow: auto;
    width: 266px;
    position: absolute;
    left: 0;
    background-color: #fff;
    z-index: 100;
    top: 34px;
}

.autocomplete-result {
    list-style: none;
    text-align: left;
    padding: 4px 12px;
    cursor: pointer;
    font-size: 110%;
}

.autocomplete-result:hover,
.is-active {
    background-color: #ddd;
    color: black;
}
</style>