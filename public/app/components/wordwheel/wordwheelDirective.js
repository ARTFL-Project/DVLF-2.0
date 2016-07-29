(function() {
    'use strict';
    angular
        .module('DVLF')
        .directive('wordwheel', wordwheel);

    function wordwheel($http, $timeout, $log) {
        return {
            templateUrl: "app/components/wordwheel/wordwheel.html",
            link: function(scope, el, attrs) {
                var currentIndex = scope.Main.wordwheel.indexOf(attrs.head);
                scope.wordwheel = scope.Main.wordwheel.slice(currentIndex-200, currentIndex+200);
                $timeout(function() {
                    var offset = angular.element('#wordwheel a.active').offset().top - 500;
                    angular.element('#wordwheel').animate({scrollTop: offset}, 0);
                })
            }
        }
    }
})();
