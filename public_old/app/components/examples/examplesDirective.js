(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('examples', examples);

    function examples($http, $timeout, $log, $location, $routeParams) {
        return {
            templateUrl: "app/components/examples/examples.html",
            link: function(scope, element, attrs) {
                scope.voted = {};
                scope.vote = function(id, vote) {
                    if (id in scope.voted) {
                        alert("Vous avez déjà voté pour cet exemple.")
                    } else {
                        scope.voted[id] = true;
                        var headword = $routeParams.queryTerm;
                        var query = "/api/vote/" + headword + "/" + id + "/" + vote;
                        $http.get(query).then(function(response) {
                            angular.element("#" + id).text(response.data.score);
                        });
                    }
                }
                scope.addExample = function() {
                    $location.url("/exemple");
                }
            }
        }
    }
})();
