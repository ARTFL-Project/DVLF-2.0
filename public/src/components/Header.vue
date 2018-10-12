<template>
    <b-navbar variant="faded" type="light" toggleable="sm">
        <b-navbar-brand>
            <router-link class="navbar-brand" style="font-variant: small-caps; text-align: center" to="/">
                <img style="max-height: 25px; display: inline-block; padding-right: 5px; margin-top: -5px;" alt="Brand" src="../assets/images/logo.png">Dictionnaire Vivant de la Langue Française
            </router-link>
        </b-navbar-brand>
        <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>
        <b-collapse is-nav id="nav_collapse" style="font-size: 80%;">
            <b-navbar-nav class="ml-auto" style="font-size: 130%;">
                <b-nav-item href="http://artfl-project.uchicago.edu">ARTFL Project</b-nav-item>
                <b-nav-item>
                    <router-link to="/apropos">À propos du DVLF</router-link>
                </b-nav-item>
                <b-nav-item-dropdown text="Contactez-nous" right>
                    <b-dropdown-item href="#">dvlf@artfl.uchicago.edu</b-dropdown-item>
                </b-nav-item-dropdown>
            </b-navbar-nav>
        </b-collapse>
    </b-navbar>
</template>

<script>
export default {
    name: "Header"
}
</script>

<style scoped>
.navbar,
#navbar {
    border-bottom: 1px solid #eee;
    background-color: #f0f0f0;
}
.navbar-brand {
    cursor: pointer;
    color: inherit !important;
    white-space: inherit;
    margin-right: 0;
    padding: 0;
    line-height: 80%;
}
/deep/ .navbar-toggler {
    padding: 0.12rem 0.35rem;
}
/deep/ .navbar-toggler-icon {
    height: 1.2em;
    width: 1.2em;
}

.navbar-nav > li > a:hover {
    color: #000 !important;
}

.nav .dropdown-menu > li > a:hover {
    color: black !important;
}
</style>