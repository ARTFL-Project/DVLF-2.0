<template>
    <nav class="navbar">
        <div class="container">
            <div class="navbar-header">
                <router-link class="navbar-brand" style="font-size: 22px; font-variant: small-caps; text-align: center" to="/">
                    <img style="max-height: 25px; display: inline-block; padding-right: 5px; margin-top: -5px;" alt="Brand" src="../assets/images/logo.png">Dictionnaire Vivant de la Langue Française
                </router-link>
                <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar" style="border: 1px solid #000; border-radius: 0px;">
                    <span class="sr-only">Toggle navigation</span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                </button>
            </div>
            <div id="navbar" class="navbar-right collapse navbar-collapse" style="font-size: 80%;">
                <ul class="nav navbar-nav">
                    <li class="active">
                        <a href="http://artfl-project.uchicago.edu">ARTFL Project</a>
                    </li>
                    <li>
                        <router-link to="/apropos">À propos du DVLF</router-link>
                    </li>
                    <li class="dropdown">
                        <a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">Contactez-nous
                            <span class="caret"></span>
                        </a>
                        <ul class="dropdown-menu">
                            <li>
                                <a>dvlf@artfl.uchicago.edu</a>
                            </li>
                        </ul>
                    </li>
                </ul>
            </div>
            <!--/.nav-collapse -->
        </div>
    </nav>
</template>

<script>
export default {
    name: "Header"
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.navbar,
#navbar {
    border-bottom: 1px solid #eee;
    background-color: #f0f0f0;
}
.navbar-brand {
    cursor: pointer;
    color: inherit !important;
}

.navbar-nav > li > a:hover {
    color: #000 !important;
}

.nav .dropdown-menu > li > a:hover {
    color: black !important;
}
</style>