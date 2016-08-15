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
                var start = currentIndex - 500;
                if (start < 0) {
                    start = 0;
                }
                var end = currentIndex + 500;
                if (end > scope.Main.wordwheel.length - 1) {
                    end = scope.Main.wordwheel.length - 1;
                }
                scope.wordwheel = scope.Main.wordwheel.slice(start, end);
                $timeout(function() {
                    var offset = angular.element('#wordwheel a.active').offset().top - 570;
                    angular.element('#wordwheel .list-group').scrollTop(offset);
                })
            }
        }
    }
})();
