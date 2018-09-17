<template>
    <div>
        <div style="margin: 15px 0px;" v-if="results.totalDicos == 0">
            Il n'existe aucune entrée pour le terme
            <b>{{ currentTerm }}</b>
            <button class="btn btn-default" style="margin-left: 10px;" type="button" @click="define()">
                <span class="glyphicon glyphicon-plus add-button"></span>
                Ajoutez votre définition
            </button>
        </div>
        <div v-if="results.totalEntries < 2">
            <div v-if="fuzzyResults.length > 0">
                Voici des termes similaires présents dans l'index:<br>
                <div style="display: inline-block" v-for="(word, index) in fuzzyResults" :key="index">
                    <router-link :to="`/mot/${word.word}`">{{word.word}}&nbsp;&nbsp;</router-link>
                </div>
            </div>
            <div v-if="fuzzyResults.length == 0">
                Aucun terme similaire dans l'index
            </div>
        </div>
        <transition-group appear name="fade">
            <div class="panel panel-default dico" style="margin-top: 15px; animation-duration: 0.4s" v-for="(dico, dicoIndex) in results.data" :key="dicoIndex">
                <div class="row dico-title">
                    <div class="col-xs-1 arrow" @click="toggleEntry(dicoIndex)">
                        <span class="glyphicon glyphicon-menu-down" v-show="dico.show"></span>
                        <span class="glyphicon glyphicon-menu-right" v-show="!dico.show"></span>
                    </div>
                    <div class="col-xs-11 dico-label" @click="toggleEntry(dicoIndex)">
                        <span class="hidden-xs hidden-sm">
                            {{ dico.label }}
                        </span>
                        <span class="hidden-md hidden-lg">
                            {{ dico.shortLabel }}
                        </span>
                        <span class="badge pull-right">{{ dico.contentObj.length }}</span>
                    </div>
                </div>
                <div class="dico-entries" :class="{showdico: dico.show}">
                    <div class="entry" style="padding-right: 10px" v-for="(entry, index) in dico.contentObj" :key="index" :class="{'tlfi': dico.name === 'tlfi'}">
                        <div v-html="entry.content"></div>
                        <div class="example-source" v-if="entry.date">Contribuée le {{ entry.date}}</div>
                        <hr v-if="!index == dico.contentObj.length-1" style="width: 50%; margin-top: 15px; margin-bottom: 5px; border-color: #ddd;">
                    </div>
                    <div class="tlfi-link" v-if="dico.name === 'tlfi'">
                        <a :href="`http://www.cnrtl.fr/definition/${encodeURIComponent(currentTerm)}`" target="_blank">=> Voir la définition complète au CNRTL</a>
                    </div>
                </div>
                <div style="margin: 10px 10px 10px 10px" v-if="dico.name === 'userSubmit'">
                    <button class="btn btn-default" type="button" @click="define()">
                        <span class="glyphicon glyphicon-plus add-button"></span>
                        Ajoutez votre définition
                    </button>
                </div>
            </div>
        </transition-group>
    </div>
</template>

<script>
export default {
    name: "DictionaryEntries",
    props: {
        results: Object,
        fuzzyResults: Array
    },
    data: function() {
        return {
            currentTerm: this.$route.params.queryTerm
        }
    },
    methods: {
        define: function() {
            this.$router.push(`/definition/${this.currentTerm}`)
        },
        toggleEntry(dicoIndex) {
            let targetElement = event.srcElement
            if (
                targetElement.classList.contains("hidden-xs") ||
                targetElement.classList.contains("badge") ||
                targetElement.classList.contains("glyphicon")
            ) {
                targetElement = targetElement.parentNode
            }
            let dicoNode = targetElement.parentNode.parentNode.parentNode
            let entry = dicoNode.querySelectorAll(".dico-entries")[dicoIndex]
            let arrows = targetElement.parentNode.parentNode.querySelectorAll(
                ".arrow .glyphicon"
            )
            if (!entry.classList.contains("showdico")) {
                entry.classList.add("showdico")
                arrows[1].style.display = "none"
                arrows[0].style.display = "initial"
            } else {
                entry.classList.remove("showdico")
                arrows[0].style.display = "none"
                arrows[1].style.display = "initial"
            }
        }
    }
}
</script>

<style scoped>
.xml-prononciation::before {
    content: "(";
}

.xml-prononciation::after {
    content: ")\A";
}

.xml-nature {
    font-style: italic;
}

.xml-indent,
.xml-variante {
    display: block;
}

.xml-variante {
    padding-top: 10px;
    padding-bottom: 10px;
    text-indent: -1.3em;
    padding-left: 1.3em;
}

.xml-variante::before {
    counter-increment: section;
    content: counter(section) ")\00a0";
    font-weight: 700;
}

:not(.xml-rubrique) + .xml-indent {
    padding-top: 10px;
}

.xml-indent {
    padding-left: 1.3em;
}

.xml-cit {
    padding-left: 2.3em;
    display: block;
    text-indent: -1.3em;
}

.xml-indent > .xml-cit {
    padding-left: 1em;
}

.xml-cit::before {
    content: "\2012\00a0\00ab\00a0";
}

.xml-cit::after {
    content: "\00a0\00bb\00a0(" attr(aut) "\00a0" attr(ref) ")";
    font-variant: small-caps;
}

.xml-rubrique {
    display: block;
    margin-top: 20px;
}

.xml-rubrique::before {
    content: attr(nom);
    font-variant: small-caps;
    font-weight: 700;
}

.xml-corps + .xml-rubrique {
    margin-top: 10px;
}

.dico {
    overflow: hidden;
}

.showdico {
    display: block !important;
}

.dico-title {
    text-align: center;
    font-size: 110%;
    font-weight: 700;
    margin-top: 0px;
    padding: 5px;
    background-color: #f0f0f0;
    border-bottom: 1px solid #eee;
}

.dico-label,
.arrow {
    cursor: pointer;
}

.dico-entries {
    text-align: justify;
    margin-top: 10px;
    line-height: 150%;
    display: none;
    transition: all 200ms;
}

.entry {
    /*max-height: 400px;*/
    overflow-y: auto;
    padding: 0 10px 10px 10px;
    font-size: inherit !important;
}

.entry.tlfi {
    max-height: 145px;
    overflow: hidden;
}

.arrow .glyphicon {
    vertical-align: middle;
    font-size: 80%;
}

.tlfi-link {
    float: left;
    padding: 10px 10px 10px 10px;
}

.tlfi-link a {
    cursor: pointer;
}
.badge {
    background-color: #fff;
    color: rgb(21, 95, 131);
    font-family: initial;
}
</style>


