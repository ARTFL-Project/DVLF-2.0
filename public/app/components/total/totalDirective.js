(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('total', total);

    function total($http, $rootScope, $routeParams, totalResults) {
        return {
            templateUrl: "app/components/total/total.html",
            link: function(scope, el, attrs) {
                if ($routeParams.queryTerm !== 'undefined') {
                    scope.queryTerm = $routeParams.queryTerm;
                } else {
                    scope.queryTerm = totalResults.queryTerm;
                }
                $rootScope.$on('resultsUpdate', function(e, args) {
                    scope.queryTerm = args.queryTerm;
                    scope.totalDicos = args.totalDicos;
                    scope.totalEntries = args.totalEntries;
                })
            }
        }
    }
})();
