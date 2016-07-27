(function() {
    "use strict";

    angular
        .module("DVLF", ['ngRoute', 'ngTouch', 'ngSanitize']);

    angular
        .module("DVLF")
        .controller('MainController', MainController);

		function MainController($routeParams, $route, $log, $http, $location) {
            var vm = this;

            vm.queryTerm = ""

            vm.apropos = false;
            vm.viewAPropos = function() {
                if (vm.apropos) {
                    vm.apropos = false;
                } else {
                    vm.apropos = true;
                }
            }
        }
})();
