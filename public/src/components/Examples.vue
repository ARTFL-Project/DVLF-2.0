<template>
    <div id="examples" class="panel panel-default">
        <div id="examples-title">
            Exemples d'utilisation
        </div>
        <div class="panel panel-default example" v-for="(example, index) in examples" :key="index">
            <div class="example-score">
                <span class="glyphicon glyphicon-circle-arrow-up up" @click="vote(example.id, 'up')"></span>
                <span :id="example.id">{{ example.score }}</span>
                <span class="glyphicon glyphicon-circle-arrow-down down" @click="vote(example.id, 'down')"></span>
            </div>
            <div class="example-content" v-html="example.content"></div>
            <div class="example-source" v-if="example.source">Source trouvée sur:
                <a :href="example.link" target="_blank" v-if="example.link != 'http://'">{{ example.source }}</a>
                <span v-if="example.link == 'http://'">{{ example.source }}</span>
            </div>
            <div class="example-source" v-if="example.date">Contribué le {{ example.date }}</div>
        </div>
        <button class="btn btn-default" style="margin: 0 10px 10px 10px" type="button" @click="addExample()">
            <span class="glyphicon glyphicon-plus add-button"></span>
            Ajoutez votre exemple
        </button>
    </div>
</template>

<script>
export default {
    name: "Examples",
    props: {
        examples: Array
    },
    data() {
        return {
            headword: this.$route.params.queryTerm,
            voted: {}
        }
    },
    methods: {
        vote(id, vote) {
            if (id in this.voted) {
                alert("Vous avez déjà voté pour cet exemple.")
            } else {
                this.voted[id] = true
                let query = `${this.$globalConfig.apiServer}/api/vote/${
                    this.headword
                }/${id}/${vote}`
                this.$http.get(query).then(function(response) {
                    document.getElementById(id).textContent =
                        response.data.score
                })
            }
        },
        addExample() {
            this.$router.push(`/exemple/${this.headword}`)
        }
    }
}
</script>

<style scoped>
#examples-title {
    text-align: center;
    font-size: 110%;
    font-weight: 700;
    margin-top: 0px;
    padding: 5px;
    background-color: #f0f0f0;
    border-bottom: 1px solid #eee;
    margin-bottom: 10px;
}

.example {
    margin: 10px;
    padding: 10px;
}

.example-score {
    display: inline;
    font-family: serif;
    font-size: 110%;
    padding-right: 5px;
}

.example-score .glyphicon {
    color: #a8a8a8;
    cursor: pointer;
    transition: all 300ms;
}

.example-score .glyphicon:hover {
    color: black;
}

.example-score .up {
    padding-right: 5px;
}

.example-score .down {
    padding-left: 5px;
}

.example-content {
    display: inline;
    padding-left: 5px;
    text-align: justify;
}

.example-source {
    display: inline-block;
    width: 100%;
    text-align: right;
    font-size: 85%;
}
</style>
