const gulp = require('gulp');
const uglify = require('gulp-uglify');
const concat = require('gulp-concat');
const hash = require('gulp-hash-filename');
const babel = require('gulp-babel');
const inject = require('gulp-inject');
var clean = require('gulp-clean');

gulp.task('clean bundle', function () {
  return gulp.src('./public/*.min.js', {read: false})
    .pipe(clean());
});

gulp.task('bundle js', function () {
  return gulp.src(
    [
      'node_modules/babel-polyfill/dist/polyfill.js',
      'src/index.js',
    ]
  ) 
    .pipe(concat('bundle.js'))
    .pipe(babel({
      presets: ['@babel/preset-env'],
    }))
    .pipe(uglify())
    .pipe(hash({
      "format": "{name}.{hash}.min{ext}"
    }))
    .pipe(gulp.dest('./public'));
});

gulp.task('inject to html', function() {
  var target = gulp.src('./public/index.html');
  const sources = gulp.src(['./public/*.min.js'], { read: false });

  return target.pipe(inject(sources, { 
      relative: true,
      transform: function(filePath) {
        return `<script type="application/javascript" src="${filePath}"></script>`;
      }
   }))
    .pipe(gulp.dest('./public'));
});

gulp.task('default', gulp.series(
  'clean bundle', 
  'bundle js', 
  'inject to html'
));