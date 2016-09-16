(function() {
    "use strict";

    angular
        .module("DVLF", ['ngRoute', 'ngTouch', 'ngSanitize', 'ngAnimate', 'angucomplete-alt', 'sticky', 'vcRecaptcha']);

    angular
        .module("DVLF")
        .controller('MainController', MainController);

    function MainController($scope, $routeParams, $route, $log, $http, $location, totalResults) {
        var vm = this;

        if (angular.equals({}, $routeParams)) {
            vm.atHome = true;
        } else {
            vm.atHome = false;
        }
        $scope.$watch(function() {
            if (angular.equals({}, $routeParams) && $location.path() != "/apropos" && $location.path() != "/synonyme" && $location.path() != "/antonyme" && $location.path() != "/definition" && $location.path() != "/exemple") {
                return true;
            } else {
                return false;
            }
        }, function(atHome) {
            if (atHome) {
                vm.atHome = true;
            } else {
                vm.atHome = false;
            }
        });

        vm.apropos = false;
        vm.viewAPropos = function() {
            $location.path("/apropos")
        }

        vm.search = function(word) {
            console.log(angular.element("#search_value").val())
            word = angular.element("#search_value").val();
            $location.path("/mot/" + word.trim());
        }

        $scope.autocomplete = function(query) {
            if (typeof(query) !== 'undefined') {
                $location.path("/mot/" + query.title.trim());
            }
        }
    }
})();
