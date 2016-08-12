(function() {
    "use strict";

    angular
        .module("DVLF")
        .controller('ResultsController', ResultsController);

		function ResultsController($scope, $log, $http, $routeParams) {
            var vm = this;
            vm.currentTerm = $routeParams.queryTerm;

        }
})();
