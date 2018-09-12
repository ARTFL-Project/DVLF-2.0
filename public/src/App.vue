<template>
    <div id="app">
        <Header/>
        <div class="container-fluid">
            <Search/>
            <router-view :key="$route.fullPath"></router-view>
        </div>
    </div>
</template>

<script>
import Header from "./components/Header.vue"
import Search from "./components/Search.vue"
import Results from "./components/Results.vue"

export default {
    name: "app",
    components: {
        Header,
        Search,
        Results
    }
}
</script>

<style>
body {
    font-size: 160%;
    counter-reset: section;
    /* Set the section counter to 0 */
}
/* Button colors */
.btn-primary {
    background: rgb(21, 95, 131);
    transition: background-color 200ms;
    border-color: transparent;
}

.btn-primary:hover,
.btn-primary:focus {
    background: rgba(21, 95, 131, 0.8);
    border-color: transparent;
    outline: 0;
}

.btn-primary.active,
.open > .dropdown-toggle.btn-primary,
#collocation_table th {
    background: rgb(21, 67, 95);
    border-color: transparent;
    color: #fff;
}

.btn-primary[disabled] {
    background-color: rgb(91, 150, 173);
    border-color: transparent;
}

/* Links */
a {
    color: rgb(21, 95, 131);
    font-weight: 400;
    transition: all 200ms;
}

a:hover {
    color: rgb(128, 0, 0) !important;
    text-decoration: none;
}

.ui-menu-item .highlight {
    color: rgb(174, 89, 80);
}

#waiting,
#waiting-inner,
#waiting-inner-inner {
    border-color: rgba(174, 89, 80, 0.5) rgba(143, 57, 49, 0.7)
        rgba(174, 89, 80, 0.5) rgba(174, 89, 80, 0.5);
}

.filter-list {
    color: rgb(109, 50, 39);
    font-weight: 400;
}

.navbar-right a:hover {
    color: #fff !important;
}

.btn-default {
    background-color: #fbfbfb;
}

.list-group-item.active {
    background-color: rgb(21, 95, 131);
    border-color: rgb(21, 95, 131);
}

.hidden {
    opacity: 0;
}
#app {
    color: rgb(61, 61, 61);
    margin: 0 auto;
    font-family: "Amiri", serif;
    font-weight: 300;
}
.badge {
    padding-top: 5px;
}

.a-propos {
    text-align: justify;
}

a.btn:hover {
    color: #fff !important;
}

.btn,
.panel,
.panel-heading,
.list-group-item:first-child,
.list-group-item:last-child,
.popover,
.progress,
.form-control,
.tooltip-inner,
.modal-content {
    border-radius: 0px !important;
}

.icon-bar {
    border-radius: 0px !important;
    border: 1px solid #000;
}

a.list-group-item:focus,
a.list-group-item:hover,
button.list-group-item:focus,
button.list-group-item:hover {
    background-color: #f0f0f0;
}

.hide-results {
    visibility: hidden;
}

.add-button {
    color: rgb(21, 95, 131);
    vertical-align: top;
}

.glyphicon.spinning {
    vertical-align: sub;
    animation: spin 1s infinite linear;
    -webkit-animation: spin2 1s infinite linear;
}

.nym-title,
.index,
#time-series-title {
    background-color: #f0f0f0;
    border-bottom: 1px solid #eee;
    text-align: center;
    margin-top: 0;
    padding: 7px;
    font-weight: 700;
}

.highlight {
    color: red;
}
</style>
