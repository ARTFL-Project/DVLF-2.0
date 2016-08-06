(function() {
    "use strict";

    angular
        .module("DVLF", ['ngRoute', 'ngTouch', 'ngSanitize', 'bootstrap3-typeahead']);

    angular
        .module("DVLF")
        .controller('MainController', MainController);

		function MainController($routeParams, $route, $log, $http, $location) {
            var vm = this;

            vm.queryTerm = ""

            $http.get('/api/wordwheel').then(function(response) {
                vm.wordwheel = response.data;
            });

            vm.apropos = false;
            vm.viewAPropos = function() {
                $location.path("/apropos")
            }

            vm.search = function(word) {
                $location.path("/mot/" + word.trim());
            }
        }
})();
