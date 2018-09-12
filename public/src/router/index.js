import Vue from 'vue';
import VueRouter from "vue-router";
import Results from "../components/Results";
import APropos from "../components/APropos";
import NewDefinition from "../components/NewDefinition";
import NewExample from "../components/NewExample";
import NewSynAnto from "../components/NewSynAnto";

Vue.use(VueRouter);

export default new VueRouter({
    mode: "history",
    routes: [{
            path: "/mot/:queryTerm",
            name: "Results",
            component: Results
        }, {
            path: "/apropos",
            name: "APropos",
            component: APropos
        },
        {
            path: "/definition/:term?",
            name: "NewDefinition",
            component: NewDefinition
        },
        {
            path: "/exemple/:term?",
            name: "NewExample",
            component: NewExample
        },
        {
            path: "/synonyme/:term?",
            name: "NewSynAnto",
            component: NewSynAnto,
            props: {
                typeOfNym: "synonyme"
            }
        },
        {
            path: "/antonyme/:term?",
            name: "NewSynAnt",
            component: NewSynAnto,
            props: {
                typeOfNym: "antonyme"
            }
        }
    ],
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition;
        } else {
            return {
                x: 0,
                y: 0
            };
        }
    }
});