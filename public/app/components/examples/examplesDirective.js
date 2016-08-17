(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('examples', examples);

    function examples($http, $timeout, $log, $location, $routeParams) {
        return {
            templateUrl: "app/components/examples/examples.html",
            link: function(scope, element, attrs) {
                scope.vote = function(id, vote) {
                    console.log(id)
                    var headword = $routeParams.queryTerm;
                    var query = "/api/vote/" + headword + "/" + id + "/" + vote;
                    $http.get(query).then(function(response) {
                        angular.element("#" + id).text(response.data.score);
                    });
                }
                scope.addExample = function() {
                    $location.url("/exemple");
                }
            }
        }
    }
})();
