(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('wordwheel', wordwheel);

    function wordwheel($http, $timeout, $rootScope, $anchorScroll) {
        return {
            templateUrl: "app/components/wordwheel/wordwheel.html",
            link: function(scope, el, attrs) {
                scope.$on("resultsUpdate", function(e, args) {
                    scope.currentTerm = args.queryTerm;
                    $http.get('/api/wordwheel?headword=' + args.queryTerm).then(function(response) {
                        scope.wordwheel = response.data;
                        $timeout(function() {
                            if (typeof(angular.element('#wordwheel a.active').offset()) != "undefined") {
                                var offset = angular.element('#wordwheel a.active').offset().top - 570;
                                if (offset != 0) {
                                    angular.element('#wordwheel .list-group').scrollTop(offset);
                                }
                            }
                        });
                    });
                })
            }
        }
    }
})();
