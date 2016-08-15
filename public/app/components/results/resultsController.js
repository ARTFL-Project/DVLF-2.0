(function() {
    "use strict";

    angular
        .module("DVLF")
        .controller('ResultsController', ResultsController);

		function ResultsController($scope, $log, $location, $routeParams) {
            var vm = this;
            vm.currentTerm = $routeParams.queryTerm;
            $scope.Main.queryTerm = vm.currentTerm;

            vm.define = function() {
                $location.path("/definition");
            }
        }
})();
