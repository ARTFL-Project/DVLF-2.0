<template>
    <div id="wordwheel" class="panel panel-default">
        <h4 class="index">Index</h4>
        <div class="list-group">
            <button type="button" id="up" class="btn btn-default" @click="addBefore()">&#8679;</button>
            <router-link :to="`/mot/${word}`" :id="word" class="list-group-item wordwheel-item" :href="`/mot/${word}`" :class="{'active': word === headword}" v-for="(word, index) in wordwheel" :key="index">
                {{ word }}
            </router-link>
            <button type="button" id="down" class="btn btn-default" @click="addAfter()">&#8679;</button>
        </div>
    </div>
</template>

<script>
import VueScrollTo from "vue-scrollto"

export default {
    name: "WordWheel",
    props: {
        headword: String
    },
    data() {
        return {
            wordwheel: [],
            firstLoad: true,
            startIndex: 0,
            endIndex: 0,
            currentFirst: null
        }
    },
    created() {
        this.fetchData()
    },
    methods: {
        fetchData() {
            let vm = this
            this.$http
                .get(
                    `${this.$globalConfig.apiServer}/api/wordwheel?headword=${
                        this.headword
                    }`,
                    {
                        headers: {
                            "Access-Control-Allow-Origin": "*",
                            "Content-Type": "application/json"
                        }
                    }
                )
                .then(function(response) {
                    vm.wordwheel = response.data.words
                    vm.startIndex = response.data.startIndex
                    vm.endIndex = response.data.endIndex
                    vm.$nextTick(function() {
                        let options = {
                            container: "#wordwheel .list-group",
                            duration: 5,
                            offset: -366
                        }
                        let element = document
                            .getElementById("wordwheel")
                            .querySelector(".active")
                        VueScrollTo.scrollTo(element, options)
                    })
                })
                .catch(error => {
                    this.error = error.toString()
                    console.log(error)
                })
        },
        addBefore() {
            this.currentFirst = this.wordwheel[0]
            var vm = this
            this.$http
                .get(
                    `${this.$globalConfig.apiServer}/api/wordwheel?startIndex=${
                        this.startIndex
                    }`
                )
                .then(function(response) {
                    vm.wordwheel.unshift(...response.data.words)
                    vm.startIndex = response.data.startIndex
                    vm.endIndex = response.data.endIndex
                    vm.$nextTick(function() {
                        let options = {
                            container: "#wordwheel .list-group",
                            duration: 10,
                            offset: -24
                        }
                        let element = document.getElementById(vm.currentFirst)
                        VueScrollTo.scrollTo(element, options)
                    })
                })
                .catch(error => {
                    this.error = error.toString()
                    console.log(error)
                })
        },
        addAfter() {
            var vm = this
            this.$http
                .get(
                    `${this.$globalConfig.apiServer}/api/wordwheel?endIndex=${
                        this.endIndex
                    }`
                )
                .then(function(response) {
                    vm.wordwheel.push(...response.data.words)
                    vm.startIndex = response.data.startIndex
                    vm.endIndex = response.data.endIndex
                })
                .catch(error => {
                    this.error = error.toString()
                    console.log(error)
                })
        }
    }
}
</script>

<style scoped>
.list-group {
    margin-top: -10px;
    max-height: 800px;
    overflow-y: scroll;
}
#up {
    padding: 4px 10px 0px;
}
#up,
#down {
    display: inline-block;
    width: 100%;
    border-width: 0px;

    background-color: #f8f8f8;
}
#down {
    padding: 0px 10px;
    transform: rotate(180deg);
}
</style>
